// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package events

import (
	"encoding/json"

	log "github.com/golang/glog"

	"ipe/subscription"
)

// SubscribeData data for Subscribe event
type SubscribeData struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth,omitempty"`
	ChannelData string `json:"channel_data,omitempty"`
}

// Subscribe event
// {
//     "event": "pusher:subscribe",
//     "data": {
//         "channel": "the channel",
//         "auth": "the auth",
//         "channelData": "extra data"
//     }
// }
type Subscribe struct {
	Event string        `json:"event"`
	Data  SubscribeData `json:"data"`
}

// NewSubscribe Create a new subscribe event with the specified channel and data
func NewSubscribe(channel, auth, channelData string) Subscribe {
	data := SubscribeData{Channel: channel, Auth: auth, ChannelData: channelData}
	return Subscribe{Event: "pusher:subscribe", Data: data}
}

// UnsubscribeData event data
type UnsubscribeData struct {
	Channel string `json:"channel"`
}

// Unsubscribe event
// {
//     "event": "pusher:unsubscribe",
//     "data": {
//         "channel": "The channel"
//     }
// }
type Unsubscribe struct {
	Event string          `json:"event"`
	Data  UnsubscribeData `json:"data"`
}

// NewUnsubscribe Create a new unsubscribe event for the specified channel
func NewUnsubscribe(channel string) Unsubscribe {
	data := UnsubscribeData{Channel: channel}
	return Unsubscribe{Event: "pusher:unsubscribe", Data: data}
}

// SubscriptionSucceeded event
// {
//     "event": "pusher_internal:subscription_succeeded",
//     "channel": "the channel"
// }
type SubscriptionSucceeded struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// NewSubscriptionSucceeded Create a new subscription succeed event for the specified channel
func NewSubscriptionSucceeded(channel, data string) SubscriptionSucceeded {
	return SubscriptionSucceeded{Event: "pusher_internal:subscription_succeeded", Channel: channel, Data: data}
}

// SubscriptionSucceededPresenceData Data Subscription Succeed
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
type SubscriptionSucceededPresenceData struct {
	Ids   []string               `json:"ids"`
	Hash  map[string]interface{} `json:"hash"`
	Count int                    `json:"count"`
}

// NewSubscriptionSucceedPresenceData returns new SubscriptionSucceededPresenceData
func NewSubscriptionSucceedPresenceData(subscriptions map[string]*subscription.Subscription) SubscriptionSucceededPresenceData {
	event := SubscriptionSucceededPresenceData{}

	var (
		ids  []string
		hash = make(map[string]interface{}, len(subscriptions))
	)

	for _, s := range subscriptions {
		// Do you have any other idea?
		var js interface{}
		if err := json.Unmarshal([]byte(s.Data), &js); err != nil {
			continue
		}

		hash[s.ID] = js
		ids = append(ids, s.ID)
	}

	event.Ids = ids
	event.Hash = hash
	event.Count = len(subscriptions)

	return event
}

// Pong event
// {
//     "event": "pusher:pong",
//     "data": {}
// }
type Pong struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// NewPong Create a new pong event
func NewPong() Pong {
	return Pong{Event: "pusher:pong", Data: "{}"}
}

// Ping event
// {
//     "event": "pusher:ping",
//     "data": {}
// }
type Ping struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// NewPing Create a new ping event
func NewPing() Ping {
	return Ping{Event: "pusher:ping", Data: "{}"}
}

// Error event
// {
//     "event": "pusher:error",
//     "data": {
//         "message": "A Message",
//         "code": 4000
//     }
// }
type Error struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// NewError Create a new error event
// Pusher protocol is very strange in some parts
// It send null in some errors.
func NewError(code int, message string) Error {
	var data = struct {
		Code    *int   `json:"code"`
		Message string `json:"message"`
	}{
		Message: message,
	}

	if code == 0 {
		data.Code = nil
	} else {
		data.Code = &code
	}

	return Error{Event: "pusher:error", Data: data}
}

// ConnectionEstablished event
// {
//     "event" : "pusher:connection_established",
//     "data" : {
//       "socket_id" : "123456",
//       "activity_timeout" : 120
//     }
// }
type ConnectionEstablished struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// NewConnectionEstablished Create a new connection established event using the specified socketId
func NewConnectionEstablished(socketID string) ConnectionEstablished {
	b, err := json.Marshal(struct {
		SocketID        string `json:"socket_id"`
		ActivityTimeout int    `json:"activity_timeout"`
	}{
		SocketID: socketID, ActivityTimeout: 120,
	})

	if err != nil {
		panic("events: Could not Marshal json ConnectionEstablishedEvent")
	}

	return ConnectionEstablished{Event: "pusher:connection_established", Data: string(b)}
}

// MemberAdded event
// {
//   "event": "pusher_internal:member_added",
//   "channel": "presence-example-channel",
//   "data": String
// }
type MemberAdded struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// NewMemberAdded creates a new MemberAdded event
func NewMemberAdded(channel, data string) MemberAdded {
	return MemberAdded{Event: "pusher_internal:member_added", Channel: channel, Data: data}
}

// MemberRemoved event
// {
//   "event": "pusher_internal:member_removed",
//   "channel": "presence-example-channel",
//   "data": String
// }
type MemberRemoved struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// NewMemberRemoved returns a new MemberRemoved event
func NewMemberRemoved(channel string, userID string) MemberRemoved {
	data, err := json.Marshal(struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	})

	if err != nil {
		log.Error(err)
	}

	return MemberRemoved{Event: "pusher_internal:member_removed", Channel: channel, Data: string(data)}
}

// Raw event, usually used for client events
// {
//     "event": "client-?",
//     "channel": "The channel",
//     "data": {}
// }
type Raw struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

// Response event
type Response struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

// NewResponse The response event that is broadcasted to the client sockets
func NewResponse(name, channel string, data interface{}) Response {
	return Response{Event: name, Channel: channel, Data: data}
}
