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
