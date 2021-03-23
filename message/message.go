package message

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"io"
	"log"
	"net"
	"sync"
)

type (
	DataType string
)

const (
	DataTypeSys             DataType = "sys"
	DataTypeReady           DataType = "ready"
	DataTypeError           DataType = "error"
	DataTypeUsers           DataType = "users"
	DataTypeSignIn          DataType = "signIn"
	DataTypeSignUp          DataType = "signUp"
	DataTypeSignOut         DataType = "signOut"
	DataTypeAuthorized      DataType = "authorized"
	DataTypeUnAuthorized    DataType = "unauthorized"
	DataTypeChannelJoin     DataType = "channelJoin"
	DataTypeChannelMessage  DataType = "channelMessage"
	DataTypeChannelMessages DataType = "channelMessages"
	DataTypeChannelLeave    DataType = "channelLeave"
)

type Message struct {
	recipientsSessionUUID []string
	SUUID                 string              `json:"SUUID,omitempty"`
	Type                  DataType            `json:"type"`
	UserUUID              string              `json:"userUUID,omitempty"`
	User                  *rediscli.User      `json:"user,omitempty"`
	UserAccessKey         string              `json:"userAccessKey,omitempty"`
	Sys                   *DataSys            `json:"sys,omitempty"`
	Ready                 *DataReady          `json:"ready,omitempty"`
	Error                 *DataError          `json:"error,omitempty"`
	Users                 *DataUsers          `json:"users,omitempty"`
	SignIn                *DataSignIn         `json:"signIn,omitempty"`
	SignUp                *DataSignUp         `json:"signUp,omitempty"`
	SignOut               *DataSignOut        `json:"signOut,omitempty"`
	Authorized            *DataAuthorized     `json:"authorized,omitempty"`
	ChannelJoin           *DataChannelJoin    `json:"channelJoin,omitempty"`
	ChannelMessage        *DataChannelMessage `json:"channelMessage,omitempty"`
	ChannelLeave          *DataChannelLeave   `json:"channelLeave,omitempty"`
}

type DataAuthorized struct {
	UserUUID  string `json:"userUUID"`
	AccessKey string `json:"accessKey"`
}

type DataUnAuthorized struct {
	UserUUID  string `json:"userUUID"`
	AccessKey string `json:"accessKey"`
}

type Channel struct {
	conn     net.Conn
	userUUID string
}

var channelSessionsJoins = map[string]map[string]Channel{}
var channelSessionsSync = &sync.RWMutex{}

var sessionChannel = map[string]string{}
var sessionChannelSync = &sync.RWMutex{}

func channelSessionsAdd(conn net.Conn, channelUUID, sessionUUID, userUUID string) {
	channelSessionsSync.Lock()
	if _, ok := channelSessionsJoins[channelUUID]; !ok {
		channelSessionsJoins[channelUUID] = make(map[string]Channel, 0)
	}
	channelSessionsJoins[channelUUID][sessionUUID] = Channel{conn: conn, userUUID: userUUID}
	channelSessionsSync.Unlock()

	sessionChannelSync.Lock()
	sessionChannel[sessionUUID] = channelUUID
	sessionChannelSync.Unlock()

}

func channelSessionsRemove(sessionUUID string) {
	sessionChannelSync.RLock()
	channelUUID := sessionChannel[sessionUUID]
	sessionChannelSync.RUnlock()
	if channelUUID == "" {
		return
	}
	channelSessionsSync.Lock()
	if _, ok := channelSessionsJoins[channelUUID]; ok {
		delete(channelSessionsJoins[channelUUID], sessionUUID)
	}
	channelSessionsSync.Unlock()
}

type Write func(conn io.ReadWriter, op ws.OpCode, message *Message) error

func channelSessionsSendMessage(skipUserUUID, channelUUID string, write Write, message *Message) {

	channelSessionsSync.RLock()
	defer channelSessionsSync.RUnlock()
	for _, data := range channelSessionsJoins[channelUUID] {
		if skipUserUUID != "" && skipUserUUID == data.userUUID {
			log.Println(">>>>>>>>>>>>SKIP",skipUserUUID, fmt.Sprintf("%+v", message))
			continue
		}
		log.Println(">>>>>>>>>>>>SEND",skipUserUUID, fmt.Sprintf("%+v", message))
		if err := write(data.conn, ws.OpText, message); err != nil {
			log.Println(err)
		}
	}

}
