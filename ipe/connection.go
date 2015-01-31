// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// An User Connection
type Connection struct {
	SocketID string
	Socket   *websocket.Conn
}

// Create a new Subscriber
func NewConnection(socketID string, s *websocket.Conn) *Connection {
	log.Infof("Creating a new Subscriber %+v", socketID)

	return &Connection{SocketID: socketID, Socket: s}
}

// Publish the message to websocket atached to this client
func (conn *Connection) Publish(m interface{}) {
	go func() {
		if err := conn.Socket.WriteJSON(m); err != nil {
			log.Errorf("Error publishing message to connection %+v, %s", conn, err)
		}
	}()
}
