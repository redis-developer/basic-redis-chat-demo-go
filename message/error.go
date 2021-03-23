package message

import "errors"

type IError interface {
	Error() (uint32, error)
}
type Error struct {
	code uint32
	err  error
}

func (err Error) Error() (uint32, error) {
	return err.code, err.err
}

func newError(errCode uint32, err error) IError {
	return Error{code: errCode, err: err}
}

const (
	errCodeOK uint32 = iota
	errCodeJSUnmarshal
	errCodeWSRead
	errCodeWSWrite
	errCodeSignIn
	errCodeSignUp
	errCodeSignOut
	errCodeRedisChannelMessage
	errCodeRedisChannelUsers
	errCodeRedisGetSessionUUID
	errCodeRedisGetUserByUUID
	errCodeRedisChannelJoin
	errCodeUserSetOnline
)

var (
	errWSRead        = errors.New("could not read websocket connection")
	errSignIn        = errors.New("could not signIn")
	errUserSetOnline = errors.New("could not set user status as OnLine")
)
