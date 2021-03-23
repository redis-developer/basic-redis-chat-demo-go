package rediscli

import (
	"fmt"
)

func redisKeyUser(userUUID string) string {
	return fmt.Sprintf("user.%s", userUUID)
}

func redisKeySessionUUID(sessionUUID string) string {
	return fmt.Sprintf("session.%s", sessionUUID)
}
func redisKeyAccessKey(userUUID string) string {
	return fmt.Sprintf("access_key.%s", userUUID)
}

func redisKeyChannel(kind, senderUUID, userRecipientUID string) string {

	if userRecipientUID == "" {
		return fmt.Sprintf("channel.%s.public", kind)
	}

	userPoints := 0
	recipientPoints := 0

	for i := 0; i < 8; i++ {
		if senderUUID[i] <= userRecipientUID[i] {
			userPoints++
		} else {
			recipientPoints++
		}
	}

	if userPoints < recipientPoints {
		return fmt.Sprintf("channel.%s.%s", kind, senderUUID)
	}
	return fmt.Sprintf("channel.%s.%s", kind, userRecipientUID)

}

func redisKeyChannelUsers(senderUUID, recipientUUID string) string {
	return redisKeyChannel("users", senderUUID, recipientUUID)
}

func redisKeyChannelMessage(senderUUID, recipientUUID string) string {
	return redisKeyChannel("messages", senderUUID, recipientUUID)
}

func redisKeyChannelJoined(senderUUID, recipientUUID string) string {
	return redisKeyChannel("joined", senderUUID, recipientUUID)
}
