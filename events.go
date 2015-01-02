// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

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
	ChannelData string `json:"channelData,omitempty"`
}

type SubscribeEvent struct {
	Event string             `json:"event"`
	Data  SubscribeEventData `json:"data"`
}

// Create a new subscribe event with the specified channel and data
func NewSubscribeEvent(channel, auth, channelData string) SubscribeEvent {
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
func NewUnsubscribeEvent(channel string) UnsubscribeEvent {
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
}

// Create a new subscription succeed event for the specified channel
func NewSubscriptionSucceededEvent(channel string) SubscriptionSucceededEvent {
	return SubscriptionSucceededEvent{Event: "pusher_internal:subscription_succeeded", Channel: channel}
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
func NewPongEvent() PongEvent {
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
func NewPingEvent() PingEvent {
	return PingEvent{Event: "pusher:ping", Data: "{}"}
}

// {
//     "event": "pusher:error",
//     "data": {
//         "message": "A Message",
//         "code": 4000
//     }
// }
type ErrorEventData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type ErrorEvent struct {
	Event string         `json:"event"`
	Data  ErrorEventData `json:"data"`
}

// Create a new error event
func NewErrorEvent(code int, message string) ErrorEvent {
	data := ErrorEventData{Message: message, Code: code}
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
	Event string                         `json:"event"`
	Data  ConnectionEstablishedEventData `json:"data"`
}

// Create a new connection established event using the specified socketId
func NewConnectionEstablishedEvent(socketId string) ConnectionEstablishedEvent {
	data := ConnectionEstablishedEventData{SocketId: socketId, ActivityTimeout: 120}
	return ConnectionEstablishedEvent{Event: "pusher:connection_established", Data: data}
}

// {
//     "event": "client-?",
//     "channel": "The channel",
//     "data": {
//         "message": "A Message"
//     }
// }
type ClientEventData struct {
	Message string `json:"data"`
}

type ClientEvent struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    ClientEventData `json:"data"`
}

// Create a new custom client event
func NewClientEvent(name, channel, message string) ClientEvent {
	data := ClientEventData{Message: message}
	return ClientEvent{Event: "pusher:client-" + name, Channel: channel, Data: data}
}
