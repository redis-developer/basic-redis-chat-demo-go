package message

import (
	"github.com/gobwas/ws"
	"net"
)

type DataSignOut struct {
	UUID string `json:"uuid"`
}

func (p Controller) SignOut(sessionUUID string, conn net.Conn, op ws.OpCode, write Write, message *Message) IError {
	_, err := p.r.UserGet(message.UserUUID)
	if err != nil {
		return newError(errCodeSignOut, err)
	}

	p.r.UserSignOut(message.UserUUID)

	err = write(conn, op, &Message{
		Type: DataTypeSignOut,
		SignOut: &DataSignOut{
			UUID: message.SignOut.UUID,
		},
	})
	if err != nil {
		return newError(0, err)
	}

	return nil
}
