package message

type DataChannelLeave struct {
	SenderUUID    string `json:"senderUUID"`
	RecipientUUID string `json:"recipientUUID"`
}

func (p Controller) ChannelLeave(sessionUUID string, writer Write, message *Message) IError {

	channelUUID, err := p.r.ChannelLeave(message.UserUUID, message.ChannelLeave.RecipientUUID)
	if err != nil {
		return newError(0, err)
	}

	channelSessionsSendMessage("", channelUUID, writer, message)
	channelSessionsRemove(sessionUUID)

	return nil

}
