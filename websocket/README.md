## Into
Backend application base on websocket

All communications between client and server processed with websocket

Data storage made on redis
## Open WS
`const ws = WebSocket('ws://localhost:8080/ws')`
## WebSocket Events
##### Produced on websocket open 
`ws.onopen = function(event){};`
##### Produced on websocket closed
`ws.onclose = function(event){};`
##### Produced on websocket catch error 
`ws.onerror = function(event){};`
##### Produced on websocket received message from server
`ws.onmessage = function(event){};`
## Send message
#### Write a message to websocket
`ws.send(JSON.stringify(body))`

The `body` must be `JSON Object`, websocket server expected message as `JSON Object` as `string`

## Accept messages from websocket
`ws.onmessage = function(event){const body = JSON.parse(event.data);};`

The `body` contained `Stringified JSON Object`

## Kind of messages
### Ready
#### Client connected and ready to send/receive messages
> ***Response***
```
{
    "type": "ready",
    "ready": {
        "sessionUUID": "Session UUID"
    }
}
```
### Error
#### Server error response
> ***Response***
```
{
    "type": "error", 
    "error": {
        "code":1, 
        "message": "Error message"
    }
}
```
The server reply error response if request could not be processed, all kinds of messages will return error response with different `error.code` and `error.message`

This response useful for easy global error check and UI-render general error message
### User SignIn
#### Make user login
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "type": "signIn", 
    "signId": {
        "username": "User name", 
        "password": "User password"
    }
}
```
> ***Response***
```
{
    "type": "authorized": 
    "authorized": {
        "userUUID": "user UUID, 
        "accessKey": "User Access Key"
    }
}
```
User will be created if not exists 
### User SignUp
#### Create user account
> ***Request***
```
{
    "SUUID":"Session UUID",
    "type": "signUp", 
    "signUp": {
        "username": "User name", 
        "password": "User password"
    }
}
```
> ***Response***
```
{
    "type": "authorized", 
    "authorized": {
        "userUUID": 1, 
        "accessKey": "User Access Key"
    }
}
```
### User SignOut
#### Make user logout
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "type": "signOut",
    "userUUID: "User UUID", 
    "userAccessKey": "User Access Key"
}
```
> 
> ***Response***
```
{
    "type": "unauthorized", 
    "unauthorized": {
        "userUUID": "User UUID"
    }
}
```
### Users List
#### Read users list
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "type": "users",
    "userUUID: "User UUID", 
    "userAccessKey": "User Access Key"
}
```
>
> ***Response***
```
{
    "type": "users", 
    "users": {
        "total": 0,
        "received": 0,
        "users": [
            {
                "UUID": "User UUID",
                "Username": "Username",
                "Password": "Password",
                "OnLine": true
            }
        ]
    }
}
```
### Join to channel
#### Connect user to channel for read and write messages
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "userUUID": "User UUID", 
    "userAccessKey": "User Access Key", 
    "type": "channelJoin", 
    "channelJoin": {
        "recipientUUID": "User UUID"
    }
}
```
> ***Response***
```
{
    "type": "channelJoin", 
    "channelJoin": {
        "messages": [
            {
                "UUID: "Message UUID", 
                "senderUUID": "User UUID", 
                "recipientUUID": "User UUID", 
                "message": "text",
                "created_at": "time"
            }
        ], 
        "users": [
            {
                "UUID": "User UUID", 
                "username": "username"
            }
        ]
    }
}
```
When `recipientUUID` equal to `0` user will be joined to public channel
### Send message
#### Write a message from user to public or private channel
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "userUUID": "User UUID", 
    "userAccessKey": "User Access Key", 
    "type": "channelMessage", 
    "channelMessage": {
        "recipientUUID": "User UUID", 
        "message": "Message text"
    }
}
```
> ***Response***
```
{
    "type": "channelMessage", 
    "channelMessage": {
        "UUID": "Message UUID"
        "senderUUID": "User UUID",
        "recipientUUID": "User UUID", 
        "message": "Message text"
    }
}
```
When `recipientUUID` equal to `0` the message will be sent to public channel
### Read messages
#### List a messages from public or private channel
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "userUUID": "User UUID", 
    "userAccessKey": "User Access Key", 
    "type": "channelMessages", 
    "channelMessages": {
        "recipientUUID": "User UUID", 
        "offset": 1, 
        "limit": 1
    }
}
```
> ***Response***
```
{
    "type": "channelRead", 
    "channelMessages": {
        "messagesTotal": 1, 
        "messagesReceived": 1, 
        "messages": [
            {
                "senderUUID": "User UUID",
                "recipientUUID": "User UUID",
                "message": "Message text",
                "createdAt": "created time"
            }
        ]
    }
}
```
When `recipientUUID` equal to `0` the messages will be read from public channel

The `channelMessages.limit` should be less that `100`, default if `10`

The `channelMessages.offset` is the entries offset, default is `0`

All messages ordered from new to older
### Leave channel
#### Exit from a channel and stop to receive messages from it 
> ***Request***
```
{
    "SUUID": "Session UUID", 
    "userUUID": "User UUID", 
    "userAccessKey": "User Access Key", 
    "type": "channelLeave"
    "channelLeave": {
        "recipientUUID": "User UUID"
    }
}
```
> ***Response***
```
{
    "type": "channelLeave", 
    "channelLeave": {
        "recipientUUID": "User UUID"
    }
}
```
The `recipientUUID` could not be equal `0`, user can not leave public channel
