package rediscli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

const (
	keyChannelUsers           = "channelUsers"
	keyChannelMessages        = "channelMessages"
	keyChannelSenderRecipient = "channelSenderRecipient"
)

type Message struct {
	UUID          string    `json:"UUID"`
	SenderUUID    string    `json:"SenderUUID"`
	Sender        *User     `json:"Sender,omitempty"`
	RecipientUUID string    `json:"RecipientUUID"`
	Recipient     *User     `json:"Recipient,omitempty"`
	Message       string    `json:"Message"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

func (r *Redis) getKeyChannelUsers(channelUUID string) string {
	return fmt.Sprintf("%s.%s", keyChannelUsers, channelUUID)
}

func (r *Redis) getKeyChannelMessages(channelUUID string) string {
	return fmt.Sprintf("%s.%s", keyChannelMessages, channelUUID)
}

func (r *Redis) getKeyChannelSenderRecipient(senderUUID, recipientUUID string) string {
	if recipientUUID == "" {
		recipientUUID = "public"
	}
	return fmt.Sprintf("%s.%s.%s", keyChannelSenderRecipient, senderUUID, recipientUUID)
}

func (r *Redis) getChannelUUID(senderUUID, recipientUUID string) (string, error) {
	if senderUUID == "" {
		return "", errors.New("empty sender UUID")
	}
	if recipientUUID == "" {
		return "public", nil
	}
	keySenderRecipient := r.getKeyChannelSenderRecipient(senderUUID, recipientUUID)
	channelUUID, err := r.client.Get(keySenderRecipient).Result()
	if err == redis.Nil {
		keyRecipientSender := r.getKeyChannelSenderRecipient(recipientUUID, senderUUID)
		channelUUID, err = r.client.Get(keyRecipientSender).Result()
		if err == redis.Nil {
			channelUUID = uuid.NewString()
			if err := r.client.Set(keySenderRecipient, channelUUID, 0).Err(); err != nil {
				return "", err
			}
			if err := r.client.Set(keyRecipientSender, channelUUID, 0).Err(); err != nil {
				return "", err
			}
			return channelUUID, nil
		} else if err != nil {
			return "", err
		}
		return channelUUID, nil
	} else if err != nil {
		return "", err
	}
	return channelUUID, nil
}

func (r *Redis) channelJoin(channelUUID, senderUUID, recipientUUID string) error {
	key := r.getKeyChannelUsers(channelUUID)
	if err := r.client.HSet(key, senderUUID, time.Now().String()).Err(); err != nil {
		return err
	}
	if recipientUUID != "" {
		if err := r.client.HSet(key, recipientUUID, time.Now().String()).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Redis) addChannelPubSub(channelUUID string, pubSub *redis.PubSub) *ChannelPubSub {
	channelPubSub := &ChannelPubSub{
		close:  make(chan struct{}, 1),
		closed: make(chan struct{}, 1),
		pubSub: pubSub,
	}
	r.channelsPubSubSync.Lock()
	if _, ok := r.channelsPubSub[channelUUID]; !ok {
		r.channelsPubSub[channelUUID] = channelPubSub
	}
	r.channelsPubSubSync.Unlock()
	return channelPubSub
}

func (r *Redis) getChannelPubSub(channelUUID string) *ChannelPubSub {
	r.channelsPubSubSync.RLock()
	pubSub, ok := r.channelsPubSub[channelUUID]
	r.channelsPubSubSync.RUnlock()
	if !ok {
		return nil
	}
	return pubSub
}

func (r *Redis) ChannelJoin(senderUUID, recipientUUID string) (*ChannelPubSub, string, error) {

	channelUUID, err := r.getChannelUUID(senderUUID, recipientUUID)
	if err != nil {
		return nil, "", err
	}

	err = r.channelJoin(channelUUID, senderUUID, recipientUUID)
	if err != nil {
		return nil, "", err
	}
	pubSub := r.client.Subscribe(channelUUID)
	channel := r.addChannelPubSub(channelUUID, pubSub)
	return channel, channelUUID, nil
}

func (r *Redis) ChannelMessage(message *Message) (string, error) {
	channelUUID, err := r.getChannelUUID(message.SenderUUID, message.RecipientUUID)
	if err != nil {
		return "", err
	}

	buff := bytes.NewBufferString("")
	enc := json.NewEncoder(buff)
	err = enc.Encode(message)
	if err != nil {
		return "", err
	}

	err = r.client.Publish(channelUUID, buff.String()).Err()
	if err != nil {
		return "", err
	}
	key := r.getKeyChannelMessages(channelUUID)
	err = r.client.RPush(key, buff.String()).Err()
	if err != nil {
		return "", err
	}
	return channelUUID, nil
}

func (r *Redis) ChannelLeave(senderUUID, recipientUUID string) (string, error) {

	channelUUID, err := r.getChannelUUID(senderUUID, recipientUUID)
	if err != nil {
		return "", err
	}

	defer func() {
		r.channelsPubSubSync.Lock()
		delete(r.channelsPubSub, channelUUID)
		r.channelsPubSubSync.Unlock()
	}()

	r.channelsPubSubSync.RLock()
	channel, ok := r.channelsPubSub[channelUUID]
	r.channelsPubSubSync.RUnlock()

	if !ok {
		return "", errors.New("channel not found")
	}

	close(channel.close)

	timeout := time.NewTimer(time.Second * 3)
	select {
	case <-channel.closed:
		return channelUUID, nil
	case <-timeout.C:
		return "", errors.New("channel closed with timeout")
	}

}

func (r *Redis) ChannelMessagesCount(channelUUID string) (int64, error)  {
	key := r.getKeyChannelMessages(channelUUID)
	return r.client.LLen(key).Result()
}

func (r *Redis) ChannelMessages(channelUUID string, offset, limit int64) ([]*Message, error) {

	key := r.getKeyChannelMessages(channelUUID)

	log.Println("ChannelMessages", key, offset, limit)

	values, err := r.client.LRange(key, offset, limit).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]*Message, 0, len(values))
	for i := range values {
		message := &Message{}
		dec := json.NewDecoder(strings.NewReader(values[i]))
		err := dec.Decode(message)
		if err != nil {
			return nil, err
		}
		if message.SenderUUID != "" {
			user, err := r.getUserFromListByUUID(message.SenderUUID)
			if err == nil {
				message.Sender = &User{
					UUID:     user.UUID,
					Username: user.Username,
				}
			}
		}
		if message.RecipientUUID != "" {
			user, err := r.getUserFromListByUUID(message.RecipientUUID)
			if err == nil {
				message.Recipient = &User{
					UUID:     user.UUID,
					Username: user.Username,
				}
			}
		}
		messages = append(messages, message)
	}

	return messages, nil

}

func (r *Redis) ChannelUsers(channelUUID string) ([]*User, error) {

	key := r.getKeyChannelUsers(channelUUID)

	values, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)

	for userUUID, _ := range values {
		user, err := r.getUserFromListByUUID(userUUID)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}

	return users, nil

}

/*
func (r *Redis) ChannelUsersUUID(channelUUID string) ([]string, error) {

	keyChannelUsers := r.keyChannelUsers(channelUUID)
	n, err := r.client.LLen(keyChannelUsers).Result()
	if err != nil {
		return nil, err
	}
	return r.client.LRange(keyChannelUsers, 0, n).Result()
}
*/
/*
func redisChannelMessage(senderUUID, recipientUUID, textMessage string) error {
	keyChannelMessage := redisKeyChannelMessage(senderUUID, recipientUUID)
	message := Message{
		UUID:          uuid.NewString(),
		SenderUUID:    senderUUID,
		RecipientUUID: recipientUUID,
		Message:       textMessage,
		CreatedAt:     time.Now(),
	}
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}
	err = r.HSet(keyChannelMessage, message.UUID, string(data)).Err()
	if err != nil {
		return fmt.Errorf("could not HSET data: %w", err)
	}
	return nil
}

func redisChannelMessages(senderUUID, recipientUUID string) ([]Message, error) {
	keyChannelMessage := redisKeyChannelMessage(senderUUID, recipientUUID)
	values, err := r.HGetAll(keyChannelMessage).Result()
	if err != nil {
		return nil, err
	}
	messages := make([]Message, 0, len(values))
	for i := range values {
		message := Message{}
		err := json.Unmarshal([]byte(values[i]), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func redisChannelUsers(senderUUID, recipientUUID string) ([]User, error) {
	key := redisKeyChannelUsers(senderUUID, recipientUUID)
	log.Println("redisChannelUsers:", key)
	values, err := r.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("could not GHETALL: %w", err)
	}

	log.Println("values:", fmt.Sprintf("%+v", values))

	users := make([]User, 0, len(values))
	for i := range values {
		user := User{}
		err := json.Unmarshal([]byte(values[i]), &user)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal [%+v] : %w", values[i], err)
		}
		users = append(users, user)
	}

	return users, nil
}
*/
