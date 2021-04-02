# Basic Redis Chat App Demo

A basic chat application built with Golang, Websocket and Redis.

## Try it out

#### Deploy to Heroku

<p>
    <a href="https://heroku.com/deploy" target="_blank">
        <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy to Heorku" />
    </a>
</p>

#### Deploy to Google Cloud

<p>
    <a href="https://deploy.cloud.run" target="_blank">
        <img src="https://deploy.cloud.run/button.svg" alt="Run on Google Cloud" width="150px"/>
    </a>
</p>

## How it works?

![How it works](docs/screenshot001.png)

### Client Server

Communication build on websocket messages.

Client should handle messages from websocket with `onmessage(event)` event and processed it, where `event.data` 
is stringify JSON.

Send request to server available with websocket `send(data)` function, where `data` is `JSON.stringify(Object)`.

Server receive and response stringify JSON object.

Usually for not atomic message sent with `type:"example"` server respond result will be with 
same message type: `type:"example"`, for details see `Chat with websocket` block below. 

Server not guarantee that responses is ordered as requests order. 

As example:

Client requests

```
>>> Request A
>>> Request B
>>> Request C
```

May have responses in other order

```
<<< Response C
<<< Response A
<<< Response B
```

For complicated throws client should receive response for previous sent message before send next message.

As example in a pseudocode: 

```
>>> Open websocket
<<< Waiting for message type ready
>>> Send message type signIn
<<< Waiting for message type authorized
>>> Send message type channelJoin
<<< Waiting for message type channelJoin
>>> Now send multiple message to websocket and multiple receive messages 
```

#### Websocket message structure

Each of messages contain `type` property as type of message and may have payload as type name.

For the atomic message client should not expect response as request result, as example the `signOut` message is atomic.

For example: message with `type:"example"` have `example` property, it can be as `Object`, `Array`, `String`.

```
{
    type: "example",
    example: {
        foo: "bar"
    }
}
```

For handle errors on server side client should look on error message: 

```
{
    type: "error",
    error: {
        code: 0,
        error: "error message",
        payload: "additional information"
    }
}
```

### Chat with websocket

Know how to make websocket implementation for client side from server side. 

**Connect to websocket `ws[s]://<apiHost:Port>/ws`**

On first connect to websocket client receive `ready` message:

```
{
    type: "ready",
    ready: {
        sessionUUID: "123e4567-e89b-12d3-a456-426614174000" 
    }
}
```

User session saved in golang map and available before api restart, we don't want to store it in redis 
for make session flow easier.   

### User sign in

Login for chatting, if user not exist it will be created.

Send a message to websocket:
```
{
    type: "signIn",
    signIn: {
        username: "Username",
        password: "Password"
    }
}
```

Receive a message from websocket:
```
{
    type: "authorized",
    authorized: {
        userUUID: "123e4567-e89b-12d3-a456-426614174000",
        accessKey: "generated session access key"
    }
}
```

Send system message to all connected users on success:
```
{
    type: "sys",
    sys: {
        type: "signIn",
        signIn: {
            uuid: "123e4567-e89b-12d3-a456-426614174000",
            username: "Username"
        }
    }
}
```

**Redis flow**

Key for store user index by user UUID:
```
usersUUIDListIndex:<UserUUID>
```
Read user by UUID from redis KV, on exist will return user index for users list in redis:
```
GET usersUUIDListIndex:123e4567-e89b-12d3-a456-426614174000
```

Create user and save to redis list if not exists.

Key for store user index by username:
```
usersUsernameListIndex:123e4567-e89b-12d3-a456-426614174000
```

Read user by username from redis KV, on exists will return user index for users list in redis:
```
GET usersUsernameListIndex:123e4567-e89b-12d3-a456-426614174000
```

Key for store users in LList:
```
users
```

Read user from the start list by user index in list, will return json stringify on success:
```
lindex users <INDEX> 
```

User structure:

```
{
    UUID: "123e4567-e89b-12d3-a456-426614174000",
    Username: "User name",
    Password: "Password",
    AccessKey: "Access key",
    OnLine: true
}
```
Add user to the end redis list, will return number of elements in success:
```
RPUSH users <JSON Stringify user structure>
```
User index in the list is `<number of elements in list> - 1`

Save index for search user by Username:
```
SET usersUsernameListIndex:123e4567-e89b-12d3-a456-426614174000 <User Index>
```
Save index for search user by UUID: 
```
SET usersUUIDListIndex:123e4567-e89b-12d3-a456-426614174000 <User Index>
```
On error, we should remove index for search by Username:
```
DEL usersUsernameListIndex:123e4567-e89b-12d3-a456-426614174000
```
Set user online status, should expire after 60sec:
```
SETEX userStatus:123e4567-e89b-12d3-a456-426614174000 <Current time as string> 60
```

### User sign up

New user registration, will append `signIn` on successful. 

Send a message to websocket:
```
{
    type: "signUp",
    signUp: {
        username: "Username",
        password: "Password"
    }
}
```

Receive a message from websocket:
```
{
    type: "authorized",
    authorized: {
        userUUID: "123e4567-e89b-12d3-a456-426614174000",
        accessKey: "generated session access key"
    }
}
```

Each of authorized messages should contain authorized properties `userUUID` and `sessionUUID`, see details below.

Send system message to all connected users:
```
{
    type: "sys",
    sys: {
        type: "signIn",
        signIn: {
            uuid: "123e4567-e89b-12d3-a456-426614174000",
            username: "Username"
        }
    }
}
```
**Redis flow**

Look at `Sign In` redis flow, it's equals 

### Logout user

This is atomic command without message body.

Send a message to websocket:
```
{
    type: "signOut",
    userUUID: "123e4567-e89b-12d3-a456-426614174000"
}
```

Receive a message from websocket:
```
{
    type: "signOut",
    signOut: {
        uuid: "123e4567-e89b-12d3-a456-426614174000"
    }
}
```

**Redis flow**

Check if user exist.
Read user index by UUID from redis list:
```
GET usersUUIDListIndex:123e4567-e89b-12d3-a456-426614174000
```
Read user from redis list by an index:
```
LINDEX users <User Index>
```
Read user online status (it called in `UserGet` method):
```
GET userStatus:123e4567-e89b-12d3-a456-426614174000
```
On user exist and signed in, remove sign in data.

Delete user access key:
```
DEL access_key:123e4567-e89b-12d3-a456-426614174000
```
Set user offline:
```
DEL userStatus:123e4567-e89b-12d3-a456-426614174000
```
### Join to channel

After signIn/signUp client should send `channelJoin` for receive messages from specified channel.

With empty `channelJoin.recipientUUID` user will join to general channel.

For private channel set `channelJoin.recipientUUID` with valid `userUUID`.

Before channel join we should leave other channels if joined, user should have one joined channel.

Send a message to websocket:
```
{
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    sessionUUID: "123e4567-e89b-12d3-a456-426614174000",
    type: "channelJoin",
    channelJoin: {
        recipientUUID: "123e4567-e89b-12d3-a456-426614174000",
    }
}
```

Receive a message from websocket
```
{
    type: "channelJoin",
    channelJoin: {
        recipientUUID: "123e4567-e89b-12d3-a456-426614174000",
        messages: [ // array of messages in channel in desc order
            {
                UUID: "123e4567-e89b-12d3-a456-426614174000", // message UUID
                SenderUUID: "123e4567-e89b-12d3-a456-426614174000",
                Sender: {
                    UUID: "123e4567-e89b-12d3-a456-426614174000", //user UUID
                    Username: "User name"
                },
                RecipientUUID: "123e4567-e89b-12d3-a456-426614174000",
                Recipient: {
                    UUID: "123e4567-e89b-12d3-a456-426614174000", //user UUID
                    Username: "User name"
                },
                Message: "Text message",
                CreatedAt: "Message send date"
            }
        ],
        users: [ // array of joined users
            {
                UUID: "123e4567-e89b-12d3-a456-426614174000", //user UUID
	            Username: "Username",
	            OnLine: true
            } 
        ]
    }
}
```

Send system message to all connected users:

```
{
    type: "sys",
    SUUID: "123e4567-e89b-12d3-a456-426614174000",
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    user: {
        UUID: "123e4567-e89b-12d3-a456-426614174000",
        Username: "User name",
        OnLine: true
    },
    sys: {
        type: "channelJoin",
        channelJoin: {
            recipientUUID: "123e4567-e89b-12d3-a456-426614174000"
        }
    }
}
```
**Redis flow**

Leave channel if joined before channel join, see redis flow in `Leave channel` section of this README.

Read user index by UUID:
```
GET usersUUIDListIndex:123e4567-e89b-12d3-a456-426614174000
```
Read user list by user index:
```
LINDEX users <INDEX> 
```
Read channel UUID:
```
GET channelSenderRecipient:<SenderUUID>:<RecipientUUID>
```
Save joined sender to channel:
```
HSET channelUsers:<ChannelUUID> <SenderUUID> <Joined date as string>
```
Save joined recipient for private channel:
```
HSET channelUsers:<ChannelUUID> <RecipientUUID> <Joined date as string>
```
Subcribe to channel:
```
SUBSCRIBE <ChannelUUID>
```
Count message in a channel:
```
LLEN channelMessages:<ChannelUUID>
```
Read last 10 messages from a channel, if number of messages in a channel less than 10:
```
LRANGE channelMessages:<ChannelUUID> 0, 10
```
Read last 10 messages from a channel, if number of messages in a channel more than 10 `<Offset>` is `<Number of messages>-1`:
```
LRANGE channelMessages:<ChannelUUID> <Offset>, -1
```
Read channel users:
```
HGETALL channelUsers:<ChannelUUID>
```

### Leave channel

Send a message to websocket:
```
{
    SUUID: "123e4567-e89b-12d3-a456-426614174000",
    type: "channelLeave",
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    channelLeave: {
        recipientUUID: "123e4567-e89b-12d3-a456-426614174000"
    }
}
```
Receive a message from websocket, all connected users will receive it:
```
{
    SUUID: "123e4567-e89b-12d3-a456-426614174000",
    type: "channelLeave",
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    channelLeave: {
        recipientUUID: "123e4567-e89b-12d3-a456-426614174000"
    }
}
```

**Redis flow**

Get channel UUID, will return `public` on empty `recipientUUID`.

Key for private channels, first UUID is a sender(userUUID), second is recipient(userUUID):
```
channelSenderRecipient:123e4567-e89b-12d3-a456-426614174000:123e4567-e89b-12d3-a456-426614174000
```
Key for public channels:
```
channelSenderRecipient:123e4567-e89b-12d3-a456-426614174000:public
```
Read channel UUID:
```
GET channelSenderRecipient:123e4567-e89b-12d3-a456-426614174000:public
```
If channel UUID not found for sender or recipient, we should crate it for both.

Generate `channelUUID`.

Set channel UUID for sender:
```
SET channelSenderRecipient:<SenderUUID>:<RecipientUUID> <ChannelUUID>
```
Set channel UUID for recipient:
```
SET channelSenderRecipient:<RecipientUUID>:<SenderUUID> <ChannelUUID>
```
### Channel message
Send a message to websocket:
```
{
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    type: "channelMessage",
    channelMessage: {
        recipientUUID: "123e4567-e89b-12d3-a456-426614174000",
        message: "Message text"
    }
}
```
Receive a message from websocket, all joined users received it too:
```
{
    type: "channelMessage",
    channelMessage: {
        SenderUUID: "123e4567-e89b-12d3-a456-426614174000",
        Sender: {
            UUID: "123e4567-e89b-12d3-a456-426614174000",
            Username: "User name",
            OnLine: true,
        },
        RecipientUUID: "123e4567-e89b-12d3-a456-426614174000",
        Recipient: {
            UUID: "123e4567-e89b-12d3-a456-426614174000",
            Username: "User name",
            OnLine: true,
        },
        Message: "Text message",
        CreatedAt: "0000-00-00T00:00:00.000000000Z"
    }
}
```

**Redis flow**

Read channel UUID. See **Redis flow** in `Channel leave`.

Publish a message to redis PubSub:
```
PUBLISH <ChannelUUID> <Message json as string>
```
Save message in the end of redis list:
```
RPUSH channelMessages.<ChannelUUID> <Message json as string>
```
Message structure:
```
{
    UUID: "123e4567-e89b-12d3-a456-426614174000", // message UUID
    SenderUUID: "123e4567-e89b-12d3-a456-426614174000", // user UUID
    RecipientUUID: "123e4567-e89b-12d3-a456-426614174000", // user UUID, or empty for public channel
    Message: "Text message",
    CreatedAt: "0000-00-00T00:00:00.000000000Z"    
}
```
### Get users
This is atomic operation, it not expected `users` block in request.

Send a message to websocket:
```
{
    userUUID: "123e4567-e89b-12d3-a456-426614174000",
    type: "users"
}
```
Receive a message from websocket:
```
{
    type: "users",
    users: {
        total: 0,
        received: 0,
        users: [
            {
                UUID: "123e4567-e89b-12d3-a456-426614174000",
                Username: "User name",
                OnLine: true
            }
        ]
    }
}
```
**Redis flow**

Read number of users:
```
LLEN users 
```
Read all users:
```
LRANGE users 0 <Number of users>
```
## How to run it locally?

#### Copy `.env.sample` to create `.env`. And provide the values for environment variables if needed

#### Run demo

```sh
docker-compose up -d
```

Follow: http://localhost:5000