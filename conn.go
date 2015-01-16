// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// An User Connection
type Connection struct {
	SocketID string
	Socket   *websocket.Conn
}

// A Channel Subscription
type Subscription struct {
	Connection *Connection
	Id         string
	Data       string
}

// A Channel
type Channel struct {
	sync.Mutex

	CreatedAt     time.Time
	ChannelID     string
	Subscriptions map[string]*Subscription
}

// Create a new Subscription
func NewSubscription(conn *Connection, data string) *Subscription {
	return &Subscription{Connection: conn, Data: data}
}

// Return true if the channel has at least one subscriber
func (c *Channel) IsOccupied() bool {
	return c.TotalSubscriptions() > 0
}

// Check if the type of the channel is presence or is private
func (c *Channel) IsPresenceOrPrivate() bool {
	return c.IsPresence() || c.IsPrivate()
}

// Check if the type of the channel is public
func (c *Channel) IsPublic() bool {
	return !c.IsPresenceOrPrivate()
}

// Check if the type of the channel is presence
func (c *Channel) IsPresence() bool {
	return strings.HasPrefix(c.ChannelID, "presence-")
}

// Check if the type of the channel is private
func (c *Channel) IsPrivate() bool {
	return strings.HasPrefix(c.ChannelID, "private-")
}

// Get the total of subscribers
func (c *Channel) TotalSubscriptions() int {
	return len(c.Subscriptions)
}

// Get the total of users.
func (c *Channel) TotalUsers() int {
	total := make(map[string]int)

	for _, s := range c.Subscriptions {
		total[s.Id]++
	}

	return len(total)
}

// Add a new subscriber to the channel
func (c *Channel) Subscribe(a *App, conn *Connection, channelData string) error {
	log.Infof("Subscribing %s to channel %s", conn.SocketID, c.ChannelID)

	c.Lock()
	defer c.Unlock()

	subscription := NewSubscription(conn, channelData)
	c.Subscriptions[conn.SocketID] = subscription

	if c.IsPresence() {
		// User Info Data
		var info struct {
			UserID   string          `json:"user_id"`
			UserInfo json.RawMessage `json:"user_info"`
		}

		log.Infof("%+v", channelData)

		if err := json.Unmarshal([]byte(channelData), &info); err != nil {
			log.Error(err)
			return err
		}

		js, err := info.UserInfo.MarshalJSON()

		if err != nil {
			log.Error(err)
			return err
		}

		// Update the Subscription
		subscription.Id = info.UserID
		subscription.Data = string(js)

		// Publish pusher_internal:member_added
		c.PublishMemberAddedEvent(a, channelData, subscription)
		// WebHook
		a.TriggerMemberAddedHook(c, subscription)

		// pusher_internal:subscription_succeeded
		data := make(map[string]SubscriptionSucceeedEventPresenceData)
		data["presence"] = NewSubscriptionSucceedEventPresenceData(c)

		js, err = json.Marshal(data)

		if err != nil {
			log.Error(err)
			return err
		}

		conn.Publish(NewSubscriptionSucceededEvent(c.ChannelID, string(js)))
	} else {
		conn.Publish(NewSubscriptionSucceededEvent(c.ChannelID, "{}"))
	}

	// WebHook
	if c.TotalSubscriptions() == 1 {
		a.TriggerChannelOccupiedHook(c)
	}

	return nil
}

// IsSubscribed check if the user is subscribed
func (c *Channel) IsSubscribed(conn *Connection) bool {
	_, exists := c.Subscriptions[conn.SocketID]
	return exists
}

// Remove the subscriber from the channel
// It destroy the channel if the channels does not have any subscribers.
func (c *Channel) Unsubscribe(a *App, conn *Connection) error {
	log.Infof("Unsubscribing %s from channel %s", conn.SocketID, c.ChannelID)

	c.Lock()
	defer c.Unlock()

	subscription, exists := c.Subscriptions[conn.SocketID]

	if !exists {
		return errors.New("Subscription not found")
	}

	delete(c.Subscriptions, conn.SocketID)

	if c.IsPresence() {
		// Publish pusher_internal:member_removed
		c.PublishMemberRemovedEvent(a, subscription)
		// Webhook
		a.TriggerMemberRemovedHook(c, subscription)
	}

	if !c.IsOccupied() {
		// WebHook
		a.TriggerChannelVacatedHook(c)

		// Remove the empty Channel
		a.RemoveChannel(c)
	}

	return nil
}

// Create a new Channel
func NewChannel(channelID string) *Channel {
	log.Infof("Creating a new channel: %s", channelID)

	return &Channel{ChannelID: channelID, CreatedAt: time.Now(), Subscriptions: make(map[string]*Subscription)}
}

// Create a new Subscriber
func NewConnection(socketID string, s *websocket.Conn) *Connection {
	log.Infof("Creating a new Subscriber %+v", socketID)

	return &Connection{SocketID: socketID, Socket: s}
}

// Publish a MemberAddedEvent to all subscriptions
func (c *Channel) PublishMemberAddedEvent(a *App, data string, subscription *Subscription) {
	for _, subs := range c.Subscriptions {
		if subs != subscription {
			subs.Connection.Publish(NewMemberAddedEvent(c.ChannelID, data))
		}
	}
}

// Publish a MemberRemovedEvent to all subscriptions
func (c *Channel) PublishMemberRemovedEvent(a *App, subscription *Subscription) {
	for _, subs := range c.Subscriptions {
		if subs != subscription {
			subs.Connection.Publish(NewMemberRemovedEvent(c.ChannelID, subscription))
		}
	}
}

// Publish messages to all Subscribers
func (c *Channel) Publish(a *App, event RawEvent, ignore string) error {
	b, err := event.Data.MarshalJSON()

	if err != nil {
		return err
	}

	var v interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	log.Infof("Publishing message %+v to channel %s", v, c.ChannelID)

	for _, subs := range c.Subscriptions {
		if subs.Connection.SocketID != ignore {
			subs.Connection.Publish(NewResponseEvent(event.Event, event.Channel, v))
		} else {
			// Webhook
			if strings.HasPrefix(event.Event, "client-") {
				a.TriggerClientEventHook(c, subs, event.Event, v)
			}
		}
	}

	return nil
}

// Publish the message to websocket atached to this client
func (conn *Connection) Publish(m interface{}) {
	go func() {
		if err := conn.Socket.WriteJSON(m); err != nil {
			log.Errorf("Error publishing message to connection %+v, %s", conn, err)
		}
	}()
}
