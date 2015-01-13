// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimiro1/ipe/utils"
	log "github.com/golang/glog"
)

// A WebHook is sent as a HTTP POST request to the url which you specify.
// The POST request payload (body) contains a JSON document, and follows the following format:
// {
//   "time_ms": 1327078148132
//   "events": [
//     { "name": "event_name", "some": "data" }
//   ]
// }
//
// Security
// Encryption
//
// You may use a HTTP or a HTTPS url for WebHooks. In most cases HTTP is sufficient, but HTTPS can be useful if your data is sensitive or if you wish to protect against replay attacks for example.
// Authentication
//
// Since anyone could in principle send WebHooks to your application, it’s important to verify that these WebHooks originated from Pusher. Valid WebHooks will therefore contain these headers which contain a HMAC signature of the WebHook payload (body):
//
//     X-Pusher-Key: A Pusher app may have multiple tokens. The oldest active token will be used, identified by this key.
//     X-Pusher-Signature: A HMAC SHA256 hex digest formed by signing the POST payload (body) with the token’s secret.
type WebHook struct {
	TimeMs int64       `json:"time_ms"`
	Events []HookEvent `json:"events"`
}

type HookEvent struct {
	Name     string      `json:"name"`
	Channel  string      `json:"channel"`
	Event    string      `json:"event,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	SocketID string      `json:"socket_id,omitempty"`
	UserId   string      `json:"user_id,omitempty"`
}

func NewChannelOcuppiedHook(channel *Channel) HookEvent {
	return HookEvent{Name: "channel_occupied", Channel: channel.ChannelID}
}

func NewChannelVacatedHook(channel *Channel) HookEvent {
	return HookEvent{Name: "channel_vacated", Channel: channel.ChannelID}
}

func NewMemberAddedHook(channel *Channel, s *Subscriber) HookEvent {
	return HookEvent{Name: "member_added", Channel: channel.ChannelID, UserId: s.Id}
}

func NewMemberRemovedHook(channel *Channel, s *Subscriber) HookEvent {
	return HookEvent{Name: "member_removed", Channel: channel.ChannelID, UserId: s.Id}
}

func NewClientHook(channel *Channel, s *Subscription, event string, data interface{}) HookEvent {
	return HookEvent{Name: "client_event", Channel: channel.ChannelID, Event: event, Data: data, SocketID: s.Subscriber.SocketID}
}

// channel_occupied
// { "name": "channel_occupied", "channel": "test_channel" }
func (a *App) TriggerChannelOccupiedHook(c *Channel) {
	event := NewChannelOcuppiedHook(c)
	triggerHook(event.Name, a, c, event)
}

// channel_vacated
// { "name": "channel_vacated", "channel": "test_channel" }
func (a *App) TriggerChannelVacatedHook(c *Channel) {
	event := NewChannelVacatedHook(c)
	triggerHook(event.Name, a, c, event)
}

// {
//   "name": "client_event",
//   "channel": "name of the channel the event was published on",
//   "event": "name of the event",
//   "data": "data associated with the event",
//   "socket_id": "socket_id of the sending socket",
//   "user_id": "user_id associated with the sending socket" # Only for presence channels
// }
func (a *App) TriggerClientEventHook(c *Channel, s *Subscription, client_event string, data interface{}) {
	event := NewClientHook(c, s, client_event, data)

	if c.IsPresence() {
		event.UserId = s.Subscriber.Id
	}

	triggerHook(event.Name, a, c, event)
}

// {
//   "name": "member_added",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *App) TriggerMemberAddedHook(c *Channel, s *Subscriber) {
	event := NewMemberAddedHook(c, s)
	triggerHook(event.Name, a, c, event)
}

// {
//   "name": "member_removed",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *App) TriggerMemberRemovedHook(c *Channel, s *Subscriber) {
	event := NewMemberRemovedHook(c, s)
	triggerHook(event.Name, a, c, event)
}

func triggerHook(name string, app *App, c *Channel, event HookEvent) {
	if !app.WebHooks {
		log.Infof("Checking webhooks enabled for app: %+v", app)
		return
	}

	go func() {
		log.Infof("Triggering %s event", name)

		hook := WebHook{TimeMs: time.Now().Unix()}

		hook.Events = append(hook.Events, event)

		var js []byte
		var err error

		js, err = json.Marshal(hook)

		if err != nil {
			log.Errorf("Error decoding json: %+v", err)
			return
		}

		var req *http.Request

		req, err = http.NewRequest("POST", app.URLWebHook, bytes.NewReader(js))

		if err != nil {
			log.Errorf("Error creating request: %+v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Pusher-Key", app.Key)
		req.Header.Set("X-Pusher-Signature", utils.HashMAC(js, []byte(app.Secret)))

		log.V(1).Infof("%+v", req.Header)
		log.V(1).Infof("%+v", string(js))

		if _, err := http.DefaultClient.Do(req); err != nil {
			log.Errorf("Error posting %s event: %+v", name, err)
		}
	}()
}
