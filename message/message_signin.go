package message

import (
	"fmt"
	"github.com/gobwas/ws"
	"log"
	"net"
	"sync"
)

type DataSignIn struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var usersConn = map[string]net.Conn{}
var usersConnSync = &sync.RWMutex{}

func (p Controller) SignIn(sessionUUID string, conn net.Conn, op ws.OpCode, write Write, message *Message) IError {

	log.Println("SignIn", sessionUUID, fmt.Sprintf("%+v", message))

	user, err := p.r.UserAuthorize(message.SignIn.Username, message.SignIn.Password)
	if err != nil {
		log.Println("SignIn", err)
		return newError(errCodeSignIn, err)
	}

	err = write(conn, op, &Message{
		Type: DataTypeAuthorized,
		Authorized: &DataAuthorized{
			UserUUID:  user.UUID,
			AccessKey: user.AccessKey,
		},
	})
	if err != nil {
		return newError(0, err)
	}

	err = p.r.UserSetOnline(user.UUID)
	if err != nil {
		log.Println(fmt.Errorf("%s:%w", errUserSetOnline, err), sessionUUID, message)
	}

	usersConnSync.Lock()
	usersConn[message.UserUUID] = conn
	for _, conn := range usersConn {
		err := write(conn, op, p.SysSignIn(user))
		if err != nil {
			log.Println(err)
		}
	}
	usersConnSync.Unlock()

	return nil
}
