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

#### Connect to websocket `ws[s]://<apiHost:Port>/ws`

On first connect to websocket client receive `ready` message:

```
{
    type: "ready",
    ready: {
        sessionUUID: "123e4567-e89b-12d3-a456-426614174000" 
    }
}
```

#### User sign in

Login for chatting, if user not exist it will be created.

Send a message to websocket
```
{
    type: "signIn",
    signIn: {
        username: "Username",
        password: "Password"
    }
}
```

Receive a message from websocket
```
{
    type: "authorized",
    authorized: {
        userUUID: "123e4567-e89b-12d3-a456-426614174000",
        accessKey: "generated session access key"
    }
}
```

Send system message to all connected users
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

#### User sign up

New user registration, will append `signIn` on successful. 

Send a message to websocket
```
{
    type: "signUp",
    signUp: {
        username: "Username",
        password: "Password"
    }
}
```

Receive a message from websocket
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

Send system message to all connected users
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

#### Join to channel

After signIn/signUp client should send `channelJoin` for receive messages from specified channel.

With empty `channelJoin.recipientUUID` user will join to general channel.

For private channel set `channelJoin.recipientUUID` with valid `userUUID`.

Send a message to websocket
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
}
```

## How to run it locally?

#### Copy `.env.sample` to create `.env`. And provide the values for environment variables if needed

#### Run demo

```sh
docker-compose up -d
```

Follow: http://localhost:5000