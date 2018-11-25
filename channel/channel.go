// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package channel

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/golang/glog"

	"ipe/connection"
	"ipe/events"
	"ipe/subscription"
	"ipe/utils"
)

type Option func(*Channel)
type ListenerFunc func(*Channel, *subscription.Subscription)
type ClientEventListenerFunc func(*Channel, *subscription.Subscription, string, interface{})

// A Channel
type Channel struct {
	sync.RWMutex

	ID            string
	subscriptions map[string]*subscription.Subscription

	createdAt time.Time

	memberAddedListeners     []ListenerFunc
	memberRemovedListeners   []ListenerFunc
	channelOccupiedListeners []ListenerFunc
	channelVacatedListeners  []ListenerFunc
	clientEventListeners     []ClientEventListenerFunc
}

// Create a new Channel
func New(channelID string, options ...Option) *Channel {
	log.Infof("Creating a new Channel: %s", channelID)

	c := &Channel{ID: channelID, createdAt: time.Now(), subscriptions: make(map[string]*subscription.Subscription)}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithMemberAddedListener(f ListenerFunc) func(*Channel) {
	return func(c *Channel) {
		c.memberAddedListeners = append(c.memberAddedListeners, f)
	}
}

func WithMemberRemovedListener(f ListenerFunc) func(*Channel) {
	return func(c *Channel) {
		c.memberRemovedListeners = append(c.memberRemovedListeners, f)
	}
}

func WithChannelOccupiedListener(f ListenerFunc) func(*Channel) {
	return func(c *Channel) {
		c.channelOccupiedListeners = append(c.channelOccupiedListeners, f)
	}
}

func WithChannelVacatedListener(f ListenerFunc) func(*Channel) {
	return func(c *Channel) {
		c.channelVacatedListeners = append(c.channelVacatedListeners, f)
	}
}

func WithClientEventListener(f ClientEventListenerFunc) func(*Channel) {
	return func(c *Channel) {
		c.clientEventListeners = append(c.clientEventListeners, f)
	}
}

// Subscriptions returns a slice of subscriptions
func (c *Channel) Subscriptions() []*subscription.Subscription {
	c.RLock()
	defer c.RUnlock()

	var subscriptions []*subscription.Subscription

	for _, sub := range c.subscriptions {
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions
}

// Return true if the Channel has at least one subscriber
func (c *Channel) IsOccupied() bool {
	return c.TotalSubscriptions() > 0
}

// Check if the type of the Channel is presence or is private
func (c *Channel) IsPresenceOrPrivate() bool {
	return c.IsPresence() || c.IsPrivate()
}

// Check if the type of the Channel is public
func (c *Channel) IsPublic() bool {
	return !c.IsPresenceOrPrivate()
}

// Check if the type of the Channel is presence
func (c *Channel) IsPresence() bool {
	return utils.IsPresenceChannel(c.ID)
}

// Check if the type of the Channel is private
func (c *Channel) IsPrivate() bool {
	return utils.IsPrivateChannel(c.ID)
}

// Get the total of subscribers
func (c *Channel) TotalSubscriptions() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.subscriptions)
}

// Get the total of users.
func (c *Channel) TotalUsers() int {
	c.RLock()
	defer c.RUnlock()

	total := make(map[string]int)

	for _, s := range c.subscriptions {
		total[s.ID]++
	}

	return len(total)
}

// Add a new subscriber to the Channel
func (c *Channel) Subscribe(conn *connection.Connection, channelData string) error {
	log.Infof("Subscribing %s to Channel %s", conn.SocketID, c.ID)

	_subscription := subscription.New(conn, channelData)
	c.Lock()
	c.subscriptions[conn.SocketID] = _subscription
	c.Unlock()

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

		c.Lock()
		// Update the Subscription
		_subscription.ID = info.UserID
		_subscription.Data = string(js)
		c.Unlock()

		// Publish pusher_internal:member_added
		c.PublishMemberAddedEvent(channelData, _subscription)

		for _, hook := range c.memberAddedListeners {
			hook(c, _subscription)
		}

		// pusher_internal:subscription_succeeded
		data := make(map[string]events.SubscriptionSucceededPresenceData)
		data["presence"] = events.NewSubscriptionSucceedPresenceData(c.subscriptions)

		js, err = json.Marshal(data)

		if err != nil {
			log.Error(err)
			return err
		}

		conn.Publish(events.NewSubscriptionSucceeded(c.ID, string(js)))
	} else {
		conn.Publish(events.NewSubscriptionSucceeded(c.ID, "{}"))
	}

	if c.TotalSubscriptions() == 1 {
		for _, hook := range c.channelOccupiedListeners {
			hook(c, _subscription)
		}
	}

	return nil
}

// IsSubscribed check if the user is subscribed
func (c *Channel) IsSubscribed(conn *connection.Connection) bool {
	c.RLock()
	defer c.RUnlock()

	_, exists := c.subscriptions[conn.SocketID]
	return exists
}

// Remove the subscriber from the Channel
// It destroy the Channel if the channels does not have any subscribers.
func (c *Channel) Unsubscribe(conn *connection.Connection) error {
	log.Infof("unsubscribe %s from Channel %s", conn.SocketID, c.ID)

	c.RLock()
	_subscription, exists := c.subscriptions[conn.SocketID]
	c.RUnlock()

	if !exists {
		return errors.New("_subscription not found")
	}

	c.Lock()
	delete(c.subscriptions, conn.SocketID)
	c.Unlock()

	if c.IsPresence() {
		// Publish pusher_internal:member_removed
		c.PublishMemberRemovedEvent(_subscription)

		for _, hook := range c.memberRemovedListeners {
			hook(c, _subscription)
		}
	}

	if !c.IsOccupied() {
		for _, hook := range c.channelVacatedListeners {
			hook(c, _subscription)
		}
	}

	return nil
}

// Publish a MemberAddedEvent to all subscriptions
func (c *Channel) PublishMemberAddedEvent(data string, subscription *subscription.Subscription) {
	c.RLock()
	defer c.RUnlock()

	for _, subs := range c.subscriptions {
		if subs != subscription {
			subs.Connection.Publish(events.NewMemberAdded(c.ID, data))
		}
	}
}

// Publish a MemberRemovedEvent to all subscriptions
func (c *Channel) PublishMemberRemovedEvent(subscription *subscription.Subscription) {
	c.RLock()
	defer c.RUnlock()

	for _, subs := range c.subscriptions {
		if subs != subscription {
			subs.Connection.Publish(events.NewMemberRemoved(c.ID, subscription.ID))
		}
	}
}

// Publish messages to all Subscribers
func (c *Channel) Publish(event events.Raw, ignore string) error {
	c.RLock()
	defer c.RUnlock()

	b, err := event.Data.MarshalJSON()

	if err != nil {
		return err
	}

	var v interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	log.Infof("Publishing message %+v to Channel %s", v, c.ID)

	for _, subs := range c.subscriptions {
		if subs.Connection.SocketID != ignore {
			subs.Connection.Publish(events.NewResponse(event.Event, event.Channel, v))
		} else {
			if utils.IsClientEvent(event.Event) {
				for _, hook := range c.clientEventListeners {
					hook(c, subs, event.Event, v)
				}
			}
		}
	}

	return nil
}
