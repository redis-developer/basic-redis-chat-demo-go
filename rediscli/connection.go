package rediscli

import (
	"fmt"
	"time"
)

const (
	keyUserSession = "userSession"
)
func (r *Redis) getKeyUserSession(userSessionUUID string) string {
	return fmt.Sprintf("%s.%s", keyUserSession, userSessionUUID)
}

func (r *Redis) AddConnection(userSessionUUID string) error {
	key := r.getKeyUserSession(userSessionUUID)
	return r.client.Set(key, time.Now().String(), time.Hour).Err()
}

func (r *Redis) DelConnection(userSessionUUID string) error {
	key := r.getKeyUserSession(userSessionUUID)
	return r.client.Del(key).Err()
}
