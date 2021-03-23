package websocket

import "errors"

type IError interface {
	Error() (uint32, error)
}

const (
	errCode uint32 = iota
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
)

var (
	errWSRead = errors.New("could not read websocket connection")
	errSignIn = errors.New("could not signIn")
)
