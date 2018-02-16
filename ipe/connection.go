// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"time"

	log "github.com/golang/glog"
)

// socket interface to write to the client
type socket interface {
	WriteJSON(interface{}) error
}

// An User Connection
type connection struct {
	SocketID  string
	Socket    socket
	CreatedAt time.Time
	Pinger *time.Ticker
	PingFrequency time.Duration
	PongTimer *time.Timer
	PongTimeout time.Duration
	PongReceived bool
}

// Create a new Subscriber
func newConnection(socketID string, s socket) *connection {
	log.Infof("Creating a new Subscriber %+v", socketID)

	return &connection{SocketID: socketID, Socket: s, CreatedAt: time.Now(), PingFrequency: 120, PongTimeout: 10}
}

// Publish the message to websocket atached to this client
func (conn *connection) Publish(m interface{}) {
	conn.Socket.WriteJSON(m)
}
