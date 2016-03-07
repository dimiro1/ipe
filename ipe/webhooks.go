// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

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
//     X-Pusher-Key: The App Key.
//     X-Pusher-Signature: A HMAC SHA256 hex digest formed by signing the POST payload (body) with the token’s secret.
type webHook struct {
	TimeMs int64       `json:"time_ms"`
	Events []hookEvent `json:"events"`
}

type hookEvent struct {
	Name     string      `json:"name"`
	Channel  string      `json:"channel"`
	Event    string      `json:"event,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	SocketID string      `json:"socket_id,omitempty"`
	UserID   string      `json:"user_id,omitempty"`
}

func newChannelOcuppiedHook(channel *channel) hookEvent {
	return hookEvent{Name: "channel_occupied", Channel: channel.ChannelID}
}

func newChannelVacatedHook(channel *channel) hookEvent {
	return hookEvent{Name: "channel_vacated", Channel: channel.ChannelID}
}

func newMemberAddedHook(channel *channel, s *subscription) hookEvent {
	return hookEvent{Name: "member_added", Channel: channel.ChannelID, UserID: s.ID}
}

func newMemberRemovedHook(channel *channel, s *subscription) hookEvent {
	return hookEvent{Name: "member_removed", Channel: channel.ChannelID, UserID: s.ID}
}

func newClientHook(channel *channel, s *subscription, event string, data interface{}) hookEvent {
	return hookEvent{Name: "client_event", Channel: channel.ChannelID, Event: event, Data: data, SocketID: s.Connection.SocketID}
}

// channel_occupied
// { "name": "channel_occupied", "channel": "test_channel" }
func (a *app) TriggerChannelOccupiedHook(c *channel) {
	event := newChannelOcuppiedHook(c)
	triggerHook(event.Name, a, c, event)
}

// channel_vacated
// { "name": "channel_vacated", "channel": "test_channel" }
func (a *app) TriggerChannelVacatedHook(c *channel) {
	event := newChannelVacatedHook(c)
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
func (a *app) TriggerClientEventHook(c *channel, s *subscription, clientEvent string, data interface{}) {
	event := newClientHook(c, s, clientEvent, data)

	if c.IsPresence() {
		event.UserID = s.ID
	}

	triggerHook(event.Name, a, c, event)
}

// {
//   "name": "member_added",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *app) TriggerMemberAddedHook(c *channel, s *subscription) {
	event := newMemberAddedHook(c, s)
	triggerHook(event.Name, a, c, event)
}

// {
//   "name": "member_removed",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *app) TriggerMemberRemovedHook(c *channel, s *subscription) {
	event := newMemberRemovedHook(c, s)
	triggerHook(event.Name, a, c, event)
}

func triggerHook(name string, a *app, c *channel, event hookEvent) {
	if !a.WebHooks {
		log.Infof("Webhooks are not enabled for app: %s", a.Name)
		return
	}

	go func() {
		log.Infof("Triggering %s event", name)

		hook := webHook{TimeMs: time.Now().Unix()}

		hook.Events = append(hook.Events, event)

		var js []byte
		var err error

		js, err = json.Marshal(hook)

		if err != nil {
			log.Errorf("Error decoding json: %+v", err)
			return
		}

		var req *http.Request

		req, err = http.NewRequest("POST", a.URLWebHook, bytes.NewReader(js))
		if err != nil {
			log.Errorf("Error creating request: %+v", err)
			return
		}

		req.Header.Set("User-Agent", "Ipe UA; (+https://github.com/dimiro1/ipe)")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Pusher-Key", a.Key)
		req.Header.Set("X-Pusher-Signature", utils.HashMAC(js, []byte(a.Secret)))

		log.V(1).Infof("%+v", req.Header)
		log.V(1).Infof("%+v", string(js))

		resp, err := http.DefaultClient.Do(req)

		// See: http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			log.Errorf("Error posting %s event: %+v", name, err)
		}
	}()
}
