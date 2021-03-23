package message

import "github.com/redis-developer/basic-redis-chat-demo-go/rediscli"

type DataSys struct {
	Type         DataType          `json:"type"`
	Message      string            `json:"message,omitempty"`
	SignIn       *DataSignIn       `json:"signIn,omitempty"`
	ChannelJoin  *DataChannelJoin  `json:"channelJoin,omitempty"`
	ChannelLeave *DataChannelLeave `json:"channelLeave,omitempty"`
}

func SysMessage(message string) *Message {
	return &Message{
		Type: DataTypeSys,
		Sys: &DataSys{
			Type:    DataTypeSys,
			Message: message,
		},
	}
}

func (p Controller) SysSignIn(user *rediscli.User) *Message {
	return &Message{
		Type: DataTypeSys,
		Sys: &DataSys{
			Type: DataTypeSignIn,
			SignIn: &DataSignIn{
				UUID:     user.UUID,
				Username: user.Username,
			},
		},
	}
}

func (p Controller) SysChannelJoin(user *rediscli.User, recipientsUUID []string) *Message {
	return &Message{
		recipientsSessionUUID: recipientsUUID,
		Type:                  DataTypeSys,
		Sys: &DataSys{
			Type: DataTypeSignIn,
			SignIn: &DataSignIn{
				UUID:     user.UUID,
				Username: user.Username,
			},
		},
	}
}

func (p Controller) SysChannelLeave(user *rediscli.User, recipientsUUID []string) *Message {
	return &Message{
		recipientsSessionUUID: recipientsUUID,
		Type:                  DataTypeSys,
		Sys: &DataSys{
			Type: DataTypeChannelLeave,
			ChannelLeave: &DataChannelLeave{
				SenderUUID: user.UUID,
			},
		},
	}
}
