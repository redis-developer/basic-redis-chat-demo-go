package message

import (
	"github.com/gobwas/ws"
	"github.com/google/uuid"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"net"
	"time"
)

type DataChannelMessage struct {
	UUID          string         `json:"UUID"`
	Sender        *rediscli.User `json:"Sender,omitempty"`
	SenderUUID    string         `json:"SenderUUID"`
	Recipient     *rediscli.User `json:"Recipient,omitempty"`
	RecipientUUID string         `json:"RecipientUUID"`
	Message       string         `json:"Message"`
	CreatedAt     time.Time      `json:"CreatedAt"`
}

func (p Controller) ChannelMessage(sessionUUID string, conn net.Conn, op ws.OpCode, writer Write, message *Message) IError {

	channelMessage := &rediscli.Message{
		UUID:          uuid.NewString(),
		SenderUUID:    message.UserUUID,
		RecipientUUID: message.ChannelMessage.RecipientUUID,
		Message:       message.ChannelMessage.Message,
		CreatedAt:     time.Now(),
	}

	channelUUID, err := p.r.ChannelMessage(channelMessage)
	if err != nil {
		return nil
	}

	channelSessionsSendMessage(message.UserUUID, channelUUID, writer, message)

	return nil
}
