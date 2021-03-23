package message

type DataError struct {
	Code  uint32 `json:"code"`
	Error string `json:"error"`
	Payload interface{} `json:"payload"`
}

func (p Controller) Error(code uint32, err error, sessionUID string, payload interface{}) *Message {
	return &Message{
		recipientsSessionUUID: []string{sessionUID},
		Type:                  DataTypeError,
		Error: &DataError{
			Code:  code,
			Error: err.Error(),
			Payload: payload,
		},
	}
}
