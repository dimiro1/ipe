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

// This mutex is used to sync the generation of the new ID
var mutex = &sync.Mutex{}

// This variable store the current user id
// Every call to newID this variable is incremented
var currentID = 0

// A subscriber
type Subscriber struct {
	Id       string
	SocketID string
	Socket   *websocket.Conn
}

// A Channel Subscription
type Subscription struct {
	Subscriber *Subscriber
	Data       string
}

// Create a new Subscription
func NewSubscription(subscriber *Subscriber, data string) *Subscription {
	return &Subscription{Subscriber: subscriber, Data: data}
}

// A channel
type Channel struct {
	sync.Mutex

	CreatedAt     time.Time
	ChannelID     string
	Subscriptions []*Subscription
}

// Return true if the channel has at least one subscriber
func (c *Channel) IsOccupied() bool {
	return c.TotalSubscriptions() > 0
}

// Check if the type of the channel is presence or is private
func (c *Channel) IsPresenceOrPrivate() bool {
	return c.IsPresence() || c.IsPrivate()
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
// For now, totalUsers is equal to totalSubscribers
func (c *Channel) TotalUsers() int {
	return c.TotalSubscriptions()
}

// Add a new subscriber to the channel
func (c *Channel) Subscribe(a *App, s *Subscriber, data string) {
	log.Infof("Subscribing %s to channel %s", s.SocketID, c.ChannelID)

	c.Lock()
	c.Subscriptions = append(c.Subscriptions, NewSubscription(s, data))
	c.Unlock()

	if c.IsPresence() {
		// Publish pusher_internal:member_added - Para todos
		// WebHook
		a.TriggerMemberAddedHook(c, s)

		// pusher_internal:subscription_succeeded
		data := make(map[string]SubscriptionSucceeedEventPresenceData, 1)
		data["presence"] = NewSubscriptionSucceedEventPresenceData(c)

		js, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
		}

		if err := s.Publish(NewSubscriptionSucceededEvent(c.ChannelID, string(js))); err != nil {
			log.Error(err)
		}
	}

	// WebHook
	if c.TotalSubscriptions() == 1 {
		a.TriggerChannelOccupiedHook(c)
	}
}

// IsSubscribed check if the user is subscribed
func (c *Channel) IsSubscribed(s *Subscriber) bool {
	for _, subs := range c.Subscriptions {
		if subs.Subscriber == s {
			return true
		}
	}
	return false
}

// Remove the subscriber from the channel
// It destroy the channel if the channels does not have any subscribers.
func (c *Channel) Unsubscribe(a *App, s *Subscriber) error {
	log.Infof("Unsubscribing %s from channel %s", s.SocketID, c.ChannelID)

	c.Lock()
	defer c.Unlock()

	index := -1
	for i, subs := range c.Subscriptions {
		if subs.Subscriber == s {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("Subscription not found")
	}
	c.Subscriptions = append(c.Subscriptions[:index], c.Subscriptions[index+1:]...)

	if c.IsPresence() {
		// Publish pusher_internal:member_removed
		// Webhook
		a.TriggerMemberRemovedHook(c, s)
	}

	// WebHook
	if c.TotalSubscriptions() == 0 {
		a.TriggerChannelVacatedHook(c)
	}

	return nil
}

// Create a new Channel
func NewChannel(channelID string) *Channel {
	log.Infof("Creating a new channel: %s", channelID)

	return &Channel{ChannelID: channelID, CreatedAt: time.Now()}
}

// This function generate a sequencial ID
func newID() string {
	mutex.Lock()
	defer mutex.Unlock()

	currentID += 1

	return string(currentID)
}

// Create a new Subscriber
func NewSubscriber(socketID string, s *websocket.Conn) *Subscriber {
	id := newID()

	log.Infof("Creating a new Subscriber %+v with id %d", socketID, id)

	return &Subscriber{Id: id, SocketID: socketID, Socket: s}
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
		if subs.Subscriber.SocketID != ignore {
			js := NewResponseEvent(event.Event, event.Channel, v)

			if err := subs.Subscriber.Publish(js); err != nil {
				continue
			}
		} else {
			// Webhook
			if strings.HasPrefix(event.Event, "client-") {
				a.TriggerClientEventHook(c, subs, event.Event)
			}
		}
	}

	return nil
}

// Publish the message to websocket atached to this client
func (s *Subscriber) Publish(m interface{}) error {
	if err := s.Socket.WriteJSON(m); err != nil {
		log.Errorf("Error sending message to subscriber %+v, %s", s, err)
		return err
	}

	return nil
}
