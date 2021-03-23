package rediscli

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	keyUsers                  = "users"
	keyUserStatus             = "userStatus"
	keyUserChannels           = "userChannels"
	keyUserAccessKey          = "userAccessKey"
	keyUsersUUIDListIndex     = "usersUUIDListIndex"
	keyUsersUsernameListIndex = "usersUsernameListIndex"
)

type User struct {
	UUID        string `json:"UUID"`
	Username    string `json:"Username"`
	Password    string `json:"Password,omitempty"`
	AccessKey   string `json:"AccessKey,omitempty"`
	OnLine      bool   `json:"OnLine"`
	SessionUUID string `json:"-"`
}

func (r *Redis) getKeyUsers() string {
	return keyUsers
}

func (r *Redis) getKeyUsersUUIDListIndex(userUUID string) string {
	return fmt.Sprintf("%s.%s", keyUsersUUIDListIndex, userUUID)
}

func (r *Redis) getKeyUsersUsernameListIndex(username string) string {
	return fmt.Sprintf("%s.%x", keyUsersUsernameListIndex, md5.Sum([]byte(username)))
}

func (r *Redis) getUserIndexByUsername(username string) (int64, error) {

	log.Println("getUserIndexByUsername", username)

	key := r.getKeyUsersUsernameListIndex(username)
	value, err := r.client.Get(key).Result()
	if err != nil {
		return 0, err
	}
	index, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (r *Redis) getUserIndexByUUID(userUUID string) (int64, error) {

	log.Println("getUserIndexByUUID", userUUID)

	key := r.getKeyUsersUUIDListIndex(userUUID)
	value, err := r.client.Get(key).Result()
	if err != nil {
		return 0, fmt.Errorf("getUserIndexByUUID: %w", err)
	}
	index, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("getUserIndexByUUID: %w", err)
	}
	return index, nil
}

func (r *Redis) getKeyUserStatus(userUUID string) string {
	return fmt.Sprintf("%s.%s", keyUserStatus, userUUID)
}

func (r *Redis) getKeyUserChannels(userUUID string) string {
	return fmt.Sprintf("%s.%s", keyUserChannels, userUUID)
}

func (r *Redis) getKeyUserAccessKey(userUUID string) string {
	return fmt.Sprintf("%s.%s", keyUserAccessKey, userUUID)
}

func (r *Redis) UserAuthorize(username, password string) (*User, error) {

	log.Println("UserAuthorize", fmt.Sprintf("[%s|%s]", username, password))

	user, err := r.UserGet(username)
	if errors.Is(err, redis.Nil) {
		user, err = r.UserCreate(username, password)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	log.Println("UserAuthorize", fmt.Sprintf("%+v", user))

	if user.Password != password {
		return nil, errors.New("wrong password")
	}

	user.AccessKey, err = r.UserUpdateAccessKey(user.UUID)
	if err != nil {
		return nil, err
	}

	err = r.UserSetOnline(user.UUID)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *Redis) addUser(user *User) error {
	buff := bytes.NewBufferString("")
	enc := json.NewEncoder(buff)
	err := enc.Encode(user)
	if err != nil {
		return err
	}

	key := r.getKeyUsers()

	elements, err := r.client.RPush(key, buff.String()).Result()
	if err != nil {
		return nil
	}

	index := elements - 1
	keyUserUsernameIndex := r.getKeyUsersUsernameListIndex(user.Username)
	keyUserUUIDIndex := r.getKeyUsersUUIDListIndex(user.UUID)

	err = r.client.Set(keyUserUsernameIndex, fmt.Sprintf("%d", index), 0).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(keyUserUUIDIndex, fmt.Sprintf("%d", index), 0).Err()
	if err != nil {
		r.client.Del(keyUserUsernameIndex)
		return err
	}

	return nil
}

func (r *Redis) getUserFromList(userIndex int64) (*User, error) {

	key := r.getKeyUsers()

	value, err := r.client.LIndex(key, userIndex).Result()
	if err != nil {
		return nil, fmt.Errorf("getUserFromList[%d]: %w", userIndex, err)
	}

	user := &User{}

	dec := json.NewDecoder(strings.NewReader(value))
	err = dec.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("getUserFromList[%d]: %w", userIndex, err)
	}

	return user, nil
}

func (r *Redis) getUserFromListByUsername(username string) (*User, error) {

	log.Println("getUserFromListByUsername", username)

	userIndex, err := r.getUserIndexByUsername(username)
	if err != nil {
		return nil, err
	}

	user, err := r.getUserFromList(userIndex)
	if err != nil {
		return nil, err
	}

	user.OnLine = r.UserIsOnline(user.UUID)

	return user, nil

}

func (r *Redis) getUserFromListByUUID(userUUID string) (*User, error) {

	log.Println("getUserFromListByUUID", userUUID)

	userIndex, err := r.getUserIndexByUUID(userUUID)
	if err != nil {
		return nil, fmt.Errorf("getUserFromListByUUID[%s]: %w", userUUID, err)
	}

	user, err := r.getUserFromList(userIndex)
	if err != nil {
		return nil, fmt.Errorf("getUserFromListByUUID[%s]: %w", userUUID, err)
	}

	user.OnLine = r.UserIsOnline(user.UUID)

	return user, nil

}

func (r *Redis) UserCreate(username, password string) (*User, error) {

	log.Println("UserCreate", fmt.Sprintf("[%s|%s]", username, password))

	if user, err := r.getUserFromListByUsername(username); err == nil {
		return user, nil
	}

	user := &User{
		UUID:     uuid.NewString(),
		Username: username,
		Password: password,
	}

	if err := r.addUser(user); err != nil {
		return nil, err
	}

	return user, nil

}

func (r *Redis) UserGet(userUUID string) (*User, error) {
	log.Println("UserGet", userUUID)
	user, err := r.getUserFromListByUUID(userUUID)
	if err != nil {
		return nil, fmt.Errorf("UserGET[%s]: %w", userUUID, err)
	}
	return user, nil
}

func (r *Redis) UserAll() ([]*User, error) {

	key := r.getKeyUsers()

	items, err := r.client.LLen(key).Result()
	if err != nil {
		return nil, err
	}

	values, err := r.client.LRange(key, 0, items).Result()

	users := make([]*User, 0, items)

	for i := range values {
		user := &User{}
		dec := json.NewDecoder(strings.NewReader(values[i]))
		err = dec.Decode(user)
		if err != nil {
			return nil, fmt.Errorf("[%s]: %w", values[i], err)
		}
		users = append(users, user)
	}

	return users, nil

}

func (r *Redis) UserDeleteAccessKey(userUUID string) {
	key := r.getKeyUserAccessKey(userUUID)
	_ = r.client.Del(key)
}

func (r *Redis) UserUpdateAccessKey(userUUID string) (string, error) {
	key := r.getKeyUserAccessKey(userUUID)
	accessKey := uuid.New().String()

	err := r.client.Set(key, accessKey, 0).Err()
	if err != nil {
		return "", err
	}
	return accessKey, nil
}

func (r *Redis) UserSetOnline(userUUID string) error {
	key := r.getKeyUserStatus(userUUID)
	return r.client.Set(key, time.Now().String(), time.Minute).Err()
}

func (r *Redis) UserSetOffline(userUUID string) {
	key := r.getKeyUserStatus(userUUID)
	r.client.Del(key)
}

func (r *Redis) UserIsOnline(userUUID string) bool {
	key := r.getKeyUserStatus(userUUID)
	err := r.client.Get(key).Err()
	if err == nil {
		return true
	}
	return false
}

func (r *Redis) UserSignOut(userUUID string) {
	keyAccessKey := redisKeyAccessKey(userUUID)
	r.UserDeleteAccessKey(keyAccessKey)
	r.UserSetOffline(userUUID)
}
