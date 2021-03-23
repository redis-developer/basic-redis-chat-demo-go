package message

import (
	"github.com/gobwas/ws"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"log"
	"net"
)

type DataChannelJoin struct {
	RecipientUUID string              `json:"recipientUUID,omitempty"`
	Messages      []*rediscli.Message `json:"messages,omitempty"`
	Users         []*rediscli.User    `json:"users,omitempty"`
}

func (p Controller) ChannelJoin(sessionUUID string,conn net.Conn, op ws.OpCode, write Write,  message *Message) (*rediscli.ChannelPubSub, IError) {

	errI := p.ChannelLeave(sessionUUID, write, &Message{
		SUUID: message.SUUID,
		Type: DataTypeChannelLeave,
		UserUUID: message.UserUUID,
		ChannelLeave: &DataChannelLeave{
			RecipientUUID: message.ChannelJoin.RecipientUUID,
		},
	})
	if errI != nil {
		log.Println(errI)
	}

	channelSessionsRemove(sessionUUID)
	user,err := p.r.UserGet(message.UserUUID)
	if err != nil {
		return nil, newError(100, err)
	}

	channel, channelUUID, err := p.r.ChannelJoin(message.UserUUID, message.ChannelJoin.RecipientUUID)
	if err != nil {
		return nil, newError(101, err)
	}

	messagesLen, err := p.r.ChannelMessagesCount(channelUUID)
	if err != nil {
		return nil, newError(111,err)
	}

	var offset int64
	var limit int64 = 10

	if messagesLen > limit {
		offset = messagesLen - 10
		limit = -1
	}

	log.Println(">>>>>>>>>>>>", channelUUID, messagesLen, offset, limit)

	channelMessages, err := p.r.ChannelMessages(channelUUID, offset, limit)
	if err != nil {
		return nil, newError(102, err)
	}

	channelUsers, err := p.r.ChannelUsers(channelUUID)
	if err != nil {
		return nil, newError(103, err)
	}

	channelSessionsAdd(conn, channelUUID, sessionUUID, message.UserUUID)

	err = write(conn, op, &Message{
		Type: DataTypeChannelJoin,
		ChannelJoin: &DataChannelJoin{
			RecipientUUID: message.ChannelJoin.RecipientUUID,
			Messages:      channelMessages,
			Users:         channelUsers,
		},
	})
	if err != nil {
		return nil, newError(104, err)
	}

	channelSessionsSendMessage("", channelUUID, write, &Message{
		Type:     DataTypeSys,
		SUUID:    sessionUUID,
		UserUUID: message.UserUUID,
		User:     user,
		Sys: &DataSys{
			Type: DataTypeChannelJoin,
			ChannelJoin: &DataChannelJoin{
				RecipientUUID: message.ChannelJoin.RecipientUUID,
			},
		},
	})

	return channel, nil
}
