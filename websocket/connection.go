package websocket

import (
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

func connectionAdd(conn net.Conn, userSessionUUID string) {
	connectionsSync.Lock()
	connections[userSessionUUID] = Connection{
		userSessionID: userSessionUUID,
		conn:          conn,
	}
	connectionsSync.Unlock()
}

func connectionDel(userSessionUUID string) {
	connectionsSync.Lock()
	delete(connections, userSessionUUID)
	connectionsSync.Unlock()
}
