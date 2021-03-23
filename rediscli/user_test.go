package rediscli

import (
	"fmt"
	"log"
	"testing"
)

var testRedisInstance = NewRedis("localhost:44444", "")
/*
func TestRedis_UserStatus(t *testing.T) {

	testCases := []struct {
		userUUID string
		status   UserStatus
		expected UserStatus
	}{
		{
			userUUID: "test user set status online",
			status:   UserStatusOnline,
			expected: UserStatusOnline,
		}, {
			userUUID: "test user set status offline",
			status:   UserStatusOffline,
			expected: UserStatusOffline,
		},
	}

	for i := range testCases {

		key := testRedisInstance.keyUsersStatus()

		err := testRedisInstance.UserSetStatus(testCases[i].userUUID, testCases[i].status)
		if err != nil {
			t.Fatal(err)
		}

		value, err := testRedisInstance.client.HGet(key, testCases[i].userUUID).Result()
		if err != nil {
			t.Fatal(err)
		}

		if value != string(testCases[i].expected) {
			t.Fatalf("expected [%s], actual [%s]", testCases[i].expected, value)
		}

	}

}/*

func TestUserMap(t *testing.T) {

	user := &User{
		UUID:     "test user map",
		Username: "test user map",
	}

	if _, err := testRedisInstance.userMapGet("not exist user"); err != redis.Nil {
		t.Fatalf("expected error [%s], actual [%s]", redis.Nil, err)
	}

	err := testRedisInstance.userMapSet(user)
	if err != nil {
		t.Fatal(err)
	}

	userUUID, err := testRedisInstance.userMapGet(user.Username)
	if err != nil {
		t.Fatal(err)
	}
	if userUUID != user.UUID {
		t.Fatalf("expected user uuid [%s], actual [%s]", user.UUID, userUUID)
	}
}*/

func TestUserCreate(t *testing.T) {

	username := "test user"
	password := "test password"

	user, err := testRedisInstance.UserCreate(username, password)
	if err != nil {
		t.Fatal("userCreate", err)
	}

	user, err = testRedisInstance.UserGet(user.UUID)
	if err != nil {
		t.Fatal("userGet", err)
	}

	if user.Username != username {
		log.Fatalf("expecte username [%s], actual username [%s]", username, user.Username)
	}

	users, err := testRedisInstance.UserAll()
	if err != nil {
		log.Fatalln("usersAll", err)
	}

	for i := range users {
		log.Println(fmt.Sprintf("%+v", users[i]))
	}

}
