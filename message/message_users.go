package message

import (
	"github.com/go-redis/redis"
	"github.com/gobwas/ws"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"net"
)

type DataUsers struct {
	Total    int              `json:"total"`
	Received int              `json:"received"`
	Users    []*rediscli.User `json:"users"`
}

func (p *Controller) Users(sessionUUID string, conn net.Conn, op ws.OpCode, write Write) IError {

	values, err := p.r.UserAll()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return newError(0, err)
	}

	users := make([]*rediscli.User, 0, len(values))

	for i := range values {

		user := &rediscli.User{
			UUID:     values[i].UUID,
			Username: values[i].Username,
			OnLine:   p.r.UserIsOnline(values[i].UUID),
		}
		users = append(users, user)
	}

	err = write(conn, op, &Message{
		Type: DataTypeUsers,
		Users: &DataUsers{
			Total:    len(users),
			Received: len(users),
			Users:    users,
		},
	})
	if err != nil {
		return newError(0, err)
	}

	return nil
}
