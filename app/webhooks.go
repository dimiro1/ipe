// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"

	"ipe/channel"
	"ipe/subscription"
	"ipe/utils"
)

const maxTimeout = 3 * time.Second

// A webHook is sent as a HTTP POST request to the url which you specify.
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
// Since anyone could in principle send WebHooks to your application, it’s important to verify that these WebHooks originated from Pusher. Valid WebHooks will therefore contain these headers which contain a HMAC signature of the webHook payload (body):
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

func newChannelOcuppiedHook(channel *channel.Channel) hookEvent {
	return hookEvent{Name: "channel_occupied", Channel: channel.ID}
}

func newChannelVacatedHook(channel *channel.Channel) hookEvent {
	return hookEvent{Name: "channel_vacated", Channel: channel.ID}
}

func newMemberAddedHook(channel *channel.Channel, s *subscription.Subscription) hookEvent {
	return hookEvent{Name: "member_added", Channel: channel.ID, UserID: s.ID}
}

func newMemberRemovedHook(channel *channel.Channel, s *subscription.Subscription) hookEvent {
	return hookEvent{Name: "member_removed", Channel: channel.ID, UserID: s.ID}
}

func newClientHook(channel *channel.Channel, s *subscription.Subscription, event string, data interface{}) hookEvent {
	return hookEvent{Name: "client_event", Channel: channel.ID, Event: event, Data: data, SocketID: s.Connection.SocketID}
}

// channel_occupied
// { "name": "channel_occupied", "channel": "test_channel" }
func (a *Application) TriggerChannelOccupiedHook(c *channel.Channel) {
	event := newChannelOcuppiedHook(c)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := triggerHook(ctx, a, event); err != nil {
		log.Errorf("triggering webhook %+v", err)
	}
}

// TriggerChannelVacatedHook channel_vacated
// { "name": "channel_vacated", "channel": "test_channel" }
func (a *Application) TriggerChannelVacatedHook(c *channel.Channel) {
	event := newChannelVacatedHook(c)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := triggerHook(ctx, a, event); err != nil {
		log.Errorf("triggering webhook %+v", err)
	}
}

// TriggerClientEventHook client_events
// {
//   "name": "client_event",
//   "channel": "name of the channel the event was published on",
//   "event": "name of the event",
//   "data": "data associated with the event",
//   "socket_id": "socket_id of the sending socket",
//   "user_id": "user_id associated with the sending socket" # Only for presence channels
// }
func (a *Application) TriggerClientEventHook(c *channel.Channel, s *subscription.Subscription, clientEvent string, data interface{}) {
	event := newClientHook(c, s, clientEvent, data)

	if c.IsPresence() {
		event.UserID = s.ID
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := triggerHook(ctx, a, event); err != nil {
		log.Errorf("triggering webhook %+v", err)
	}
}

// TriggerMemberAddedHook member_added
// {
//   "name": "member_added",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *Application) TriggerMemberAddedHook(c *channel.Channel, s *subscription.Subscription) {
	event := newMemberAddedHook(c, s)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := triggerHook(ctx, a, event); err != nil {
		log.Errorf("triggering webhook %+v", err)
	}
}

// TriggerMemberRemovedHook member_removed
// {
//   "name": "member_removed",
//   "channel": "presence-your_channel_name",
//   "user_id": "a_user_id"
// }
func (a *Application) TriggerMemberRemovedHook(c *channel.Channel, s *subscription.Subscription) {
	event := newMemberRemovedHook(c, s)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := triggerHook(ctx, a, event); err != nil {
		log.Errorf("triggering webhook %+v", err)
	}
}

func triggerHook(ctx context.Context, a *Application, event hookEvent) error {
	if !a.WebHooks {
		log.Infof("webhook are not enabled for app: %s", a.Name)
		return fmt.Errorf("webhooks are not enabled for app: %s", a.Name)
	}

	done := make(chan bool)

	go func() {
		log.Infof("Triggering %s event", event.Name)

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

		req = req.WithContext(ctx)

		req.Header.Set("User-Agent", "Ipe UA; (+https://github.com/dimiro1/ipe)")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Pusher-Key", a.Key)
		req.Header.Set("X-Pusher-Signature", utils.HashMAC(js, []byte(a.Secret)))

		log.V(1).Infof("%+v", req.Header)
		log.V(1).Infof("%+v", string(js))

		resp, err := http.DefaultClient.Do(req)

		// See: http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
		if resp != nil {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Errorf("error closing response body %+v", err)
				}
			}()
		}

		if err != nil {
			log.Errorf("error posting %s event: %+v", event.Name, err)
		}

		// Successfully terminated
		done <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
