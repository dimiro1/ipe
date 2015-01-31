// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"net/http"
)

// A route
type Route struct {
	Name             string
	Method           string
	Pattern          string
	HandlerFunc      http.HandlerFunc
	RequiresRestAuth bool
}

type Routes []Route

var routes = Routes{
	Route{
		"PostEvents",
		"POST",
		"/apps/{app_id}/events",
		PostEvents,
		true,
	},
	Route{
		"GetChannels",
		"GET",
		"/apps/{app_id}/channels",
		GetChannels,
		true,
	},
	Route{
		"GetChannel",
		"GET",
		"/apps/{app_id}/channels/{channel_name}",
		GetChannel,
		true,
	},
	Route{
		"GetChannelUsers",
		"GET",
		"/apps/{app_id}/channels/{channel_name}/users",
		GetChannelUsers,
		true,
	},
	Route{
		"Websocket",
		"GET",
		"/app/{key}",
		Websocket,
		false,
	},
}
