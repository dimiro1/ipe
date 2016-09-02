// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/dimiro1/ipe/utils"
	log "github.com/golang/glog"
)

// A Channel
type channel struct {
	sync.Mutex

	CreatedAt     time.Time
	ChannelID     string
	Subscriptions map[string]*subscription
}

// Return true if the channel has at least one subscriber
func (c *channel) IsOccupied() bool {
	return c.TotalSubscriptions() > 0
}

// Check if the type of the channel is presence or is private
func (c *channel) IsPresenceOrPrivate() bool {
	return c.IsPresence() || c.IsPrivate()
}

// Check if the type of the channel is public
func (c *channel) IsPublic() bool {
	return !c.IsPresenceOrPrivate()
}

// Check if the type of the channel is presence
func (c *channel) IsPresence() bool {
	return utils.IsPresenceChannel(c.ChannelID)
}

// Check if the type of the channel is private
func (c *channel) IsPrivate() bool {
	return utils.IsPrivateChannel(c.ChannelID)
}

// Get the total of subscribers
func (c *channel) TotalSubscriptions() int {
	return len(c.Subscriptions)
}

// Get the total of users.
func (c *channel) TotalUsers() int {
	total := make(map[string]int)

	for _, s := range c.Subscriptions {
		total[s.ID]++
	}

	return len(total)
}

// Add a new subscriber to the channel
func (c *channel) Subscribe(a *app, conn *connection, channelData string) error {
	log.Infof("Subscribing %s to channel %s", conn.SocketID, c.ChannelID)

	c.Lock()
	defer c.Unlock()

	subscription := newSubscription(conn, channelData)
	c.Subscriptions[conn.SocketID] = subscription

	if !c.IsPresence() {
		conn.Publish(newSubscriptionSucceededEvent(c.ChannelID, "{}"))
		return nil
	}
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
	subscription.ID = info.UserID
	subscription.Data = string(js)

	// Publish pusher_internal:member_added
	c.PublishMemberAddedEvent(a, channelData, subscription)
	// WebHook
	a.TriggerMemberAddedHook(c, subscription)

	// pusher_internal:subscription_succeeded
	data := make(map[string]subscriptionSucceeedEventPresenceData)
	data["presence"] = newSubscriptionSucceedEventPresenceData(c)

	js, err = json.Marshal(data)

	if err != nil {
		log.Error(err)
		return err
	}

	conn.Publish(newSubscriptionSucceededEvent(c.ChannelID, string(js)))

	// WebHook
	if c.TotalSubscriptions() == 1 {
		a.TriggerChannelOccupiedHook(c)
	}

	return nil
}

// IsSubscribed check if the user is subscribed
func (c *channel) IsSubscribed(conn *connection) bool {
	_, exists := c.Subscriptions[conn.SocketID]
	return exists
}

// Remove the subscriber from the channel
// It destroy the channel if the channels does not have any subscribers.
func (c *channel) Unsubscribe(a *app, conn *connection) error {
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
func newChannel(channelID string) *channel {
	log.Infof("Creating a new channel: %s", channelID)

	return &channel{ChannelID: channelID, CreatedAt: time.Now(), Subscriptions: make(map[string]*subscription)}
}

// Publish a MemberAddedEvent to all subscriptions
func (c *channel) PublishMemberAddedEvent(a *app, data string, subscription *subscription) {
	for _, subs := range c.Subscriptions {
		if subs != subscription {
			subs.Connection.Publish(newMemberAddedEvent(c.ChannelID, data))
		}
	}
}

// Publish a MemberRemovedEvent to all subscriptions
func (c *channel) PublishMemberRemovedEvent(a *app, subscription *subscription) {
	for _, subs := range c.Subscriptions {
		if subs != subscription {
			subs.Connection.Publish(newMemberRemovedEvent(c.ChannelID, subscription))
		}
	}
}

// Publish messages to all Subscribers
func (c *channel) Publish(a *app, event rawEvent, ignore string) error {
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
			subs.Connection.Publish(newResponseEvent(event.Event, event.Channel, v))
		} else {
			// Webhook
			if utils.IsClientEvent(event.Event) {
				a.TriggerClientEventHook(c, subs, event.Event, v)
			}
		}
	}

	return nil
}
