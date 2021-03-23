package websocket

import (
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
)

var connectionsSync sync.RWMutex

type Connection struct {
	conn          io.ReadWriter
	userSessionID string
}

var connections = map[string]Connection{}

func connectionAdd(conn net.Conn) string {
	userSessionID := uuid.NewString()
	connectionsSync.Lock()
	connections[userSessionID] = Connection{
		userSessionID: userSessionID,
		conn:          conn,
	}
	connectionsSync.Unlock()
	return userSessionID
}

func connectionDel(userSessionID string) {
	connectionsSync.Lock()
	delete(connections, userSessionID)
	connectionsSync.Unlock()
}

/*
func connectionSendPublicMessage(message *Message) {
	connSync.RLock()
	for i := range connections {
		wsWrite(connections[i].conn, ws.OpText, message)
	}
	connSync.RUnlock()
}


func connectionSendMessage(message *Message, userSessions ...uuid.UUID) {
	connSync.RLock()
	for i := range userSessions {
		if c, ok := connections[userSessions[i]]; ok {
			wsWrite(c.conn, ws.OpText, message)
		}
	}
	connSync.RUnlock()
}
*/
