// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"time"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// An User Connection
type connection struct {
	SocketID  string
	Socket    *websocket.Conn
	CreatedAt time.Time
}

// Create a new Subscriber
func newConnection(socketID string, s *websocket.Conn) *connection {
	log.Infof("Creating a new Subscriber %+v", socketID)

	return &connection{SocketID: socketID, Socket: s, CreatedAt: time.Now()}
}

// Publish the message to websocket atached to this client
func (conn *connection) Publish(m interface{}) {
	go func() {
		if err := conn.Socket.WriteJSON(m); err != nil {
			log.Errorf("Error publishing message to connection %+v, %s", conn, err)
		}
	}()
}
