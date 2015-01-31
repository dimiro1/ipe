// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"

	log "github.com/golang/glog"
)

// {
//     "event": "pusher:subscribe",
//     "data": {
//         "channel": "the channel",
//         "auth": "the auth",
//         "channelData": "extra data"
//     }
// }
type SubscribeEventData struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth,omitempty"`
	ChannelData string `json:"channel_data,omitempty"`
}

type SubscribeEvent struct {
	Event string             `json:"event"`
	Data  SubscribeEventData `json:"data"`
}

// Create a new subscribe event with the specified channel and data
func newSubscribeEvent(channel, auth, channelData string) SubscribeEvent {
	data := SubscribeEventData{Channel: channel, Auth: auth, ChannelData: channelData}
	return SubscribeEvent{Event: "pusher:subscribe", Data: data}
}

type UnsubscribeEventData struct {
	Channel string `json:"channel"`
}

// {
//     "event": "pusher:unsubscribe",
//     "data": {
//         "channel": "The channel"
//     }
// }
type UnsubscribeEvent struct {
	Event string               `json:"event"`
	Data  UnsubscribeEventData `json:"data"`
}

// Create a new unsubscribe event for the specified channel
func newUnsubscribeEvent(channel string) UnsubscribeEvent {
	data := UnsubscribeEventData{Channel: channel}
	return UnsubscribeEvent{Event: "pusher:unsubscribe", Data: data}
}

// {
//     "event": "pusher_internal:subscription_succeeded",
//     "channel": "the channel"
// }
type SubscriptionSucceededEvent struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// Create a new subscription succeed event for the specified channel
func newSubscriptionSucceededEvent(channel, data string) SubscriptionSucceededEvent {
	return SubscriptionSucceededEvent{Event: "pusher_internal:subscription_succeeded", Channel: channel, Data: data}
}

// Data Subscription Succeed

// "{
//     \"presence\": {
//        \"ids\": [\"11814b369700141b222a3f3791cec2d9\",\"71dd6a29da2a4833336d2a964becf820\"],
//        \"hash\": {
//           \"11814b369700141b222a3f3791cec2d9\": {
//              \"name\":\"Phil Leggetter\",
//              \"twitter\": \"@leggetter\"
//           },
//           \"71dd6a29da2a4833336d2a964becf820\": {
//              \"name\":\"Max Williams\",
//              \"twitter\": \"@maxthelion\"
//           }
//        },
//        \"count\": 2
//     }
// }"
type SubscriptionSucceeedEventPresenceData struct {
	Ids   []string               `json:"ids"`
	Hash  map[string]interface{} `json:"hash"`
	Count int                    `json:"count"`
}

func newSubscriptionSucceedEventPresenceData(c *Channel) SubscriptionSucceeedEventPresenceData {
	event := SubscriptionSucceeedEventPresenceData{}

	var ids []string
	hash := make(map[string]interface{}, c.TotalSubscriptions())

	for _, s := range c.Subscriptions {
		// Do you have any other idea?
		var js interface{}
		json.Unmarshal([]byte(s.Data), &js)

		hash[s.Id] = js
		ids = append(ids, s.Id)
	}

	event.Ids = ids
	event.Hash = hash
	event.Count = c.TotalSubscriptions()

	return event
}

// {
//     "event": "pusher:pong",
//     "data": {}
// }
type PongEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Create a new pong event
func newPongEvent() PongEvent {
	return PongEvent{Event: "pusher:pong", Data: "{}"}
}

// {
//     "event": "pusher:ping",
//     "data": {}
// }
type PingEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Create a new ping event
func newPingEvent() PingEvent {
	return PingEvent{Event: "pusher:ping", Data: "{}"}
}

// {
//     "event": "pusher:error",
//     "data": {
//         "message": "A Message",
//         "code": 4000
//     }
// }
type ErrorEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// Create a new error event
// Pusher protocol is very strange in some parts
// It send null in some errors.
// So I created this GENERIC_ERROR thing, just to verify if the json must have null on the error code
func newErrorEvent(code int, message string) ErrorEvent {
	var data interface{}

	if code == GENERIC_ERROR {
		data = struct {
			Code    *int   `json:"code"`
			Message string `json:"message"`
		}{
			nil,
			message,
		}
	} else {
		data = struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			code,
			message,
		}
	}

	return ErrorEvent{Event: "pusher:error", Data: data}
}

// {
//     "event" : "pusher:connection_established",
//     "data" : {
//       "socket_id" : "123456",
//       "activity_timeout" : 120
//     }
// }
type ConnectionEstablishedEventData struct {
	SocketId        string `json:"socket_id"`
	ActivityTimeout int    `json:"activity_timeout"`
}

type ConnectionEstablishedEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Create a new connection established event using the specified socketId
func newConnectionEstablishedEvent(socketId string) ConnectionEstablishedEvent {
	data := ConnectionEstablishedEventData{SocketId: socketId, ActivityTimeout: 120}

	b, err := json.Marshal(data)

	if err != nil {
		panic("events: Could not Marshal json ConnectionEstablishedEvent")
	}

	return ConnectionEstablishedEvent{Event: "pusher:connection_established", Data: string(b)}
}

// {
//   "event": "pusher_internal:member_added",
//   "channel": "presence-example-channel",
//   "data": String
// }
type MemberAddedEvent struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

func newMemberAddedEvent(channel, data string) MemberAddedEvent {
	return MemberAddedEvent{Event: "pusher_internal:member_added", Channel: channel, Data: data}
}

// {
//   "event": "pusher_internal:member_removed",
//   "channel": "presence-example-channel",
//   "data": String
// }
type MemberRemovedEvent struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

func newMemberRemovedEvent(channel string, s *Subscription) MemberRemovedEvent {
	data, err := json.Marshal(struct {
		UserID string `json:"user_id"`
	}{
		UserID: s.Id,
	})

	if err != nil {
		log.Error(err)
	}

	return MemberRemovedEvent{Event: "pusher_internal:member_removed", Channel: channel, Data: string(data)}
}

// {
//     "event": "client-?",
//     "channel": "The channel",
//     "data": {}
// }
type RawEvent struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type ResponseEvent struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

// The response event that is broadcasted to the client sockets
func newResponseEvent(name, channel string, data interface{}) ResponseEvent {
	return ResponseEvent{Event: name, Channel: channel, Data: data}
}
