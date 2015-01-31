// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"errors"
	"expvar"
	"fmt"
	"sync"

	log "github.com/golang/glog"
)

// An App
type App struct {
	sync.Mutex

	Name                string
	AppID               string
	Key                 string
	Secret              string
	OnlySSL             bool
	ApplicationDisabled bool
	UserEvents          bool
	WebHooks            bool
	URLWebHook          string

	Channels    map[string]*Channel    `json:"-"`
	Connections map[string]*Connection `json:"-"`

	Stats *expvar.Map `json:"-"`
}

// Alloc memory for Connections and Channels
func (a *App) Init() {
	a.Connections = make(map[string]*Connection)
	a.Channels = make(map[string]*Channel)
	a.Stats = expvar.NewMap(fmt.Sprintf("%s (%s)", a.Name, a.AppID))
}

// Only Presence channels
func (a *App) PresenceChannels() []*Channel {
	var channels []*Channel

	for _, c := range a.Channels {
		if c.IsPresence() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Only Private channels
func (a *App) PrivateChannels() []*Channel {
	var channels []*Channel

	for _, c := range a.Channels {
		if c.IsPrivate() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Only Public channels
func (a *App) PublicChannels() []*Channel {
	var channels []*Channel

	for _, c := range a.Channels {
		if c.IsPublic() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Disconnect Socket
func (a *App) Disconnect(socketID string) {
	log.Infof("Disconnecting socket %+v", socketID)

	conn, err := a.FindConnection(socketID)

	if err != nil {
		log.Infof("Socket not found, %+v", err)
		return
	}

	// Unsubscribe from channels
	for _, c := range a.Channels {
		if c.IsSubscribed(conn) {
			c.Unsubscribe(a, conn)
		}
	}

	// Remove from app
	a.Lock()
	defer a.Unlock()

	_, exists := a.Connections[conn.SocketID]

	if !exists {
		return
	}

	delete(a.Connections, conn.SocketID)

	a.Stats.Add("TotalConnections", -1)
}

// Connect a new Subscriber
func (a *App) Connect(conn *Connection) {
	log.Infof("Adding a new Connection %s to app %s", conn.SocketID, a.Name)
	a.Lock()
	defer a.Unlock()

	a.Connections[conn.SocketID] = conn

	a.Stats.Add("TotalConnections", 1)
}

// Find a Connection on this app
func (a *App) FindConnection(socketID string) (*Connection, error) {
	conn, exists := a.Connections[socketID]

	if exists {
		return conn, nil
	}

	return nil, errors.New("Connection not found")
}

// DeleteChannel removes the channel from app
func (a *App) RemoveChannel(c *Channel) {
	log.Infof("Remove the channel %s from app %s", c.ChannelID, a.Name)
	a.Lock()
	defer a.Unlock()

	delete(a.Channels, c.ChannelID)

	if c.IsPresence() {
		a.Stats.Add("TotalPresenceChannels", -1)
	}

	if c.IsPrivate() {
		a.Stats.Add("TotalPrivateChannels", -1)
	}

	if c.IsPublic() {
		a.Stats.Add("TotalPublicChannels", -1)
	}

	a.Stats.Add("TotalChannels", -1)
}

// Add a new Channel to this APP
func (a *App) AddChannel(c *Channel) {
	log.Infof("Adding a new channel %s to app %s", c.ChannelID, a.Name)

	a.Lock()
	defer a.Unlock()

	a.Channels[c.ChannelID] = c

	if c.IsPresence() {
		a.Stats.Add("TotalPresenceChannels", 1)
	}

	if c.IsPrivate() {
		a.Stats.Add("TotalPrivateChannels", 1)
	}

	if c.IsPublic() {
		a.Stats.Add("TotalPublicChannels", 1)
	}

	a.Stats.Add("TotalChannels", 1)
}

// Returns a Channel from this app
// If not found then the channel is created and added to this app
func (a *App) FindOrCreateChannelByChannelID(n string) *Channel {
	c, err := a.FindChannelByChannelID(n)

	if err != nil {
		c = newChannel(n)
		a.AddChannel(c)
	}

	return c
}

// Find the channel by channel ID
func (a *App) FindChannelByChannelID(n string) (*Channel, error) {
	c, exists := a.Channels[n]

	if exists {
		return c, nil
	}

	return nil, errors.New("Channel does not exists")
}

func (a *App) Publish(c *Channel, event RawEvent, ignore string) error {
	a.Stats.Add("TotalUniqueMessages", 1)

	return c.Publish(a, event, ignore)
}

func (a *App) Unsubscribe(c *Channel, conn *Connection) error {
	return c.Unsubscribe(a, conn)
}

func (a *App) Subscribe(c *Channel, conn *Connection, data string) error {
	return c.Subscribe(a, conn, data)
}
