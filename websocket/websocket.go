package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/redis-developer/basic-redis-chat-demo-go/message"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func Write(conn io.ReadWriter, op ws.OpCode, message *message.Message) error {

	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("write socket message:", string(data))
	err = wsutil.WriteServerMessage(conn, op, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func NewConnection(conn net.Conn, r *rediscli.Redis, c *message.Controller, initErr chan error) {
	userSessionUUID := uuid.NewString()

	err := r.AddConnection(userSessionUUID)
	if err != nil {
		initErr <- err
		return
	}

	connectionAdd(conn, userSessionUUID)
	defer func() {
		conn.Close()
		err := r.DelConnection(userSessionUUID)
		if err != nil {
			log.Println(err)
		}
		connectionDel(userSessionUUID)

	}()

	err = Write(conn, ws.OpText, c.Ready(userSessionUUID))
	if err != nil {
		initErr <- err
		return
	}

	initErr <- nil

	for {
		msg := &message.Message{}

		if data, op, err := wsutil.ReadClientData(conn); err != nil {
			log.Println(err)
			return
			//response := makeError(errCodeWSRead, fmt.Errorf("%s: %w", errWSRead, err))
			//wsWrite(ch, conn, op, response)
		} else if err = json.Unmarshal(data, msg); err != nil {
			response := c.Error(errCodeJSUnmarshal, err, userSessionUUID, msg)
			err = Write(conn, op, response)
		} else {

			var receivedErr IError

			log.Println("Received message:", string(data))
			switch msg.Type {
			case message.DataTypeSignIn:
				receivedErr = c.SignIn(userSessionUUID, conn, op, Write, msg)
			case message.DataTypeSignUp:
				receivedErr = c.SignUp(userSessionUUID, conn, op, Write, msg)
			case message.DataTypeSignOut:
				receivedErr = c.SignOut(userSessionUUID, conn, op, Write, msg)
			case message.DataTypeUsers:
				receivedErr = c.Users(userSessionUUID, conn, op, Write)
			case message.DataTypeChannelJoin:
				channelPubSub := new(rediscli.ChannelPubSub)
				channelPubSub, receivedErr = c.ChannelJoin(userSessionUUID, conn, op, Write, msg)
				if channelPubSub != nil {
					go chatReceiver(conn, channelPubSub, r, c)
				}
			case message.DataTypeChannelMessage:
				receivedErr = c.ChannelMessage(userSessionUUID, conn, op, Write, msg)
			case message.DataTypeChannelLeave:
				receivedErr = c.ChannelLeave(userSessionUUID, Write, msg)
			default:
				err := Write(conn, op, c.Error(errCode, fmt.Errorf("unknow request data type: %s", msg.Type), msg.UserUUID, msg))
				if err != nil {
					log.Println(err)
					continue
				}
			}

			if receivedErr != nil {
				log.Println(receivedErr)
				code, err := receivedErr.Error()
				err = Write(conn, op, c.Error(code, err, userSessionUUID, string(data)))
				log.Println(receivedErr)
			}
		}
	}
}

func Handler(r *rediscli.Redis, c *message.Controller) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(request, writer)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			_,_ = fmt.Fprintf(writer, "%s", err)
			return
		}
		chInitErr := make(chan error, 1)
		go NewConnection(conn, r, c, chInitErr)
		if err = <- chInitErr; err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			_,_ = fmt.Fprintf(writer, "%s", err)
			return
		}
	}
}

func chatReceiver(conn net.Conn, channel *rediscli.ChannelPubSub, r *rediscli.Redis, c *message.Controller) {

	defer channel.Closed()
	for {
		select {
		case data := <-channel.Channel():
			msg := &message.DataChannelMessage{}
			dec := json.NewDecoder(strings.NewReader(data.Payload))
			err := dec.Decode(msg)
			if err != nil {
				log.Println(err)
			} else {

				if msg.SenderUUID != "" {
					user, err := r.UserGet(msg.SenderUUID)
					if err == nil {
						msg.Sender = &rediscli.User{
							UUID:     user.UUID,
							Username: user.Username,
						}
					}
				}

				if msg.RecipientUUID != "" {
					user, err := r.UserGet(msg.RecipientUUID)
					if err == nil {
						msg.Recipient = &rediscli.User{
							UUID:     user.UUID,
							Username: user.Username,
						}
					}
				}

				err := Write(conn, ws.OpText, &message.Message{
					Type:           message.DataTypeChannelMessage,
					ChannelMessage: msg,
				})
				if err != nil {
					log.Println(err)
				}
			}
		case <-channel.Close():
			log.Println("Close channel")
			return
		}
	}

}
