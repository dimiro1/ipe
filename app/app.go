// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	"errors"
	"expvar"
	"fmt"
	"sync"

	log "github.com/golang/glog"

	"ipe/channel"
	"ipe/connection"
	"ipe/events"
	"ipe/subscription"
)

// An App
type Application struct {
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

	channels    map[string]*channel.Channel       `json:"-"`
	connections map[string]*connection.Connection `json:"-"`

	Stats *expvar.Map `json:"-"`
}

func NewApplication(
	name,
	appID,
	key,
	secret string,
	onlySSL,
	disabled,
	userEvents,
	webHooks bool,
	webHookURL string,
) *Application {

	a := &Application{
		Name:                name,
		AppID:               appID,
		Key:                 key,
		Secret:              secret,
		OnlySSL:             onlySSL,
		ApplicationDisabled: disabled,
		UserEvents:          userEvents,
		WebHooks:            webHooks,
		URLWebHook:          webHookURL,
	}

	a.connections = make(map[string]*connection.Connection)
	a.channels = make(map[string]*channel.Channel)
	a.Stats = expvar.NewMap(fmt.Sprintf("%s (%s)", a.Name, a.AppID))

	return a
}

// Channels returns the full list of channels
func (a *Application) Channels() []*channel.Channel {
	var channels []*channel.Channel

	for _, c := range a.channels {
		channels = append(channels, c)
	}

	return channels
}

// Only Presence channels
func (a *Application) PresenceChannels() []*channel.Channel {
	var channels []*channel.Channel

	for _, c := range a.channels {
		if c.IsPresence() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Only Private channels
func (a *Application) PrivateChannels() []*channel.Channel {
	var channels []*channel.Channel

	for _, c := range a.channels {
		if c.IsPrivate() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Only Public channels
func (a *Application) PublicChannels() []*channel.Channel {
	var channels []*channel.Channel

	for _, c := range a.channels {
		if c.IsPublic() {
			channels = append(channels, c)
		}
	}

	return channels
}

// Disconnect Socket
func (a *Application) Disconnect(socketID string) {
	log.Infof("disconnecting socket %+v", socketID)

	conn, err := a.FindConnection(socketID)

	if err != nil {
		log.Infof("socket not found, %+v", err)
		return
	}

	// Unsubscribe from channels
	for _, c := range a.channels {
		if c.IsSubscribed(conn) {
			if err := c.Unsubscribe(conn); err != nil {
				log.Errorf("error while calling Channel.Unsubscribe, %+v", err)
				continue
			}
		}
	}

	// Remove from Application
	a.Lock()
	defer a.Unlock()

	_, exists := a.connections[conn.SocketID]

	if !exists {
		return
	}

	delete(a.connections, conn.SocketID)

	a.Stats.Add("TotalConnections", -1)
}

// Connect a new Subscriber
func (a *Application) Connect(conn *connection.Connection) {
	log.Infof("adding a new Connection %s to Application %s", conn.SocketID, a.Name)
	a.Lock()
	defer a.Unlock()

	a.connections[conn.SocketID] = conn

	a.Stats.Add("TotalConnections", 1)
}

// Find a Connection on this Application
func (a *Application) FindConnection(socketID string) (*connection.Connection, error) {
	conn, exists := a.connections[socketID]

	if exists {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// DeleteChannel removes the Channel from Application
func (a *Application) RemoveChannel(c *channel.Channel) {
	log.Infof("remove the Channel %s from Application %s", c.ID, a.Name)
	a.Lock()
	defer a.Unlock()

	delete(a.channels, c.ID)

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
func (a *Application) AddChannel(c *channel.Channel) {
	log.Infof("adding a new Channel %s to Application %s", c.ID, a.Name)

	a.Lock()
	defer a.Unlock()

	a.channels[c.ID] = c

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

// Returns a Channel from this Application
// If not found then the Channel is created and added to this Application
func (a *Application) FindOrCreateChannelByChannelID(n string) *channel.Channel {
	c, err := a.FindChannelByChannelID(n)

	if err != nil {
		c = channel.New(
			n,
			channel.WithChannelOccupiedListener(func(c *channel.Channel, s *subscription.Subscription) {
				a.TriggerChannelOccupiedHook(c)
			}),
			channel.WithChannelVacatedListener(func(c *channel.Channel, s *subscription.Subscription) {
				a.TriggerChannelVacatedHook(c)
			}),
			channel.WithMemberAddedListener(func(c *channel.Channel, s *subscription.Subscription) {
				a.TriggerMemberAddedHook(c, s)
			}),
			channel.WithMemberRemovedListener(func(c *channel.Channel, s *subscription.Subscription) {
				a.TriggerMemberRemovedHook(c, s)
			}),
			channel.WithClientEventListener(func(c *channel.Channel, s *subscription.Subscription, event string, data interface{}) {
				a.TriggerClientEventHook(c, s, event, data)
			}),
		)
		a.AddChannel(c)
	}

	return c
}

// Find the Channel by Channel ID
func (a *Application) FindChannelByChannelID(n string) (*channel.Channel, error) {
	c, exists := a.channels[n]

	if exists {
		return c, nil
	}

	return nil, errors.New("channel does not exists")
}

func (a *Application) Publish(c *channel.Channel, event events.Raw, ignore string) error {
	a.Stats.Add("TotalUniqueMessages", 1)

	return c.Publish(event, ignore)
}

func (a *Application) Unsubscribe(c *channel.Channel, conn *connection.Connection) error {
	err := c.Unsubscribe(conn)
	if err != nil {
		return err
	}

	if !c.IsOccupied() {
		a.RemoveChannel(c)
	}

	return nil
}

func (a *Application) Subscribe(c *channel.Channel, conn *connection.Connection, data string) error {
	return c.Subscribe(conn, data)
}
