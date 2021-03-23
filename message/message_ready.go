package message

type DataReady struct {
	SessionUUID string `json:"sessionUUID"`
}

func (p Controller) Ready(sessionUUID string) *Message {
	return &Message{
		Type: DataTypeReady,
		Ready: &DataReady{
			SessionUUID: sessionUUID,
		},
	}
}
