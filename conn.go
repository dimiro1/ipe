// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

// A subscriber
type Connection struct {
	Id       int
	SocketID string
	Data     string // Extra data attached to this subscriber
	Socket   *websocket.Conn
	Messages chan []byte
}

// A channel
type Channel struct {
	ChannelID   string
	Data        string
	Connections []*Connection
	Messages    chan []byte
}

// Return true if the channel has at least one subscriber
func (c Channel) isOccupied() bool {
	return c.totalConnections() > 0
}

// Check if the type of the channel is presence
func (c Channel) isPresence() bool {
	return strings.HasPrefix(c.ChannelID, "presence-")
}

// Check if the type of the channel is private
func (c Channel) isPrivate() bool {
	return strings.HasPrefix(c.ChannelID, "private-")
}

// Get the total of subscribers
func (c Channel) totalConnections() int {
	return len(c.Connections)
}

// Get the total of users.
// For now, totalUsers is equal to totalSubscribers
func (c Channel) totalUsers() int {
	return c.totalConnections()
}

// Add a new subscriber to the channel
func (c *Channel) Subscribe(conn *Connection) {
	c.Connections = append(c.Connections, conn)
}

// Remove the subscriber from the channel
func (c *Channel) Unsubscribe(conn *Connection) error {
	index := -1

	for i, c := range c.Connections {
		if c == conn {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("Connections not found")
	}

	c.Connections = append(c.Connections[:index], c.Connections[index+1:]...)

	// Remove Channel if necessary
	// Close sockets and Channels

	return nil
}

// Create a new Channel
func NewChannel(channelID, data string) *Channel {
	c := &Channel{ChannelID: channelID, Data: data, Messages: make(chan []byte)}
	c.Listen()

	return c
}

var mu = &sync.Mutex{}
var currentID = 0

// Generate a New ID
func newID() int {
	mu.Lock()
	defer mu.Unlock()
	currentID += 1

	return currentID
}

// Create a new Connection
func NewConnection(socketID, data string, socket *websocket.Conn) *Connection {
	id := newID()

	connection := &Connection{Id: id, SocketID: socketID, Data: data, Socket: socket, Messages: make(chan []byte)}
	connection.Listen()

	return connection
}

func (c *Channel) Listen() {
	go func() {
		select {
		case message := <-c.Messages:
			// err := c.Socket.WriteMessage(websocket.TextMessage, message)
			// log.Println(err)
			log.Println(message)
		}

	}()
}

// Listen Messages
func (c *Connection) Listen() {
	go func() {
		select {
		case message := <-c.Messages:
			// err := c.Socket.WriteMessage(websocket.TextMessage, message)
			// log.Println(err)
			log.Println(message)
		}
	}()
}
