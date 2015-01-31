// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"net/http"
)

// A route
type route struct {
	Name             string
	Method           string
	Pattern          string
	HandlerFunc      http.HandlerFunc
	RequiresRestAuth bool
}

var routes = []route{
	route{
		"PostEvents",
		"POST",
		"/apps/{app_id}/events",
		postEvents,
		true,
	},
	route{
		"GetChannels",
		"GET",
		"/apps/{app_id}/channels",
		getChannels,
		true,
	},
	route{
		"GetChannel",
		"GET",
		"/apps/{app_id}/channels/{channel_name}",
		getChannel,
		true,
	},
	route{
		"GetChannelUsers",
		"GET",
		"/apps/{app_id}/channels/{channel_name}/users",
		getChannelUsers,
		true,
	},
	route{
		"Websocket",
		"GET",
		"/app/{key}",
		wsHandler,
		false,
	},
}
