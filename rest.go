// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

// An event consists of a name and data (typically JSON) which may be sent to all subscribers to a particular channel or channels.
// This is conventionally known as triggering an event.
//
// The body should contain a Hash of parameters encoded as JSON where data parameter itself is JSON encoded.
//
// Not Implemented:
// Note that these parameters may be provided in the query string, although this is discouraged.
//
// Example:
//
// {"name":"foo","channels":["project-3"],"data":"{\"some\":\"data\"}"}
//
// Response is an empty JSON hash.
//
// POST /apps/{app_id}/events
func PostEvents(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		Data     string   `json:"data"`
		Channels []string `json:"channels,omitempty"`
		Channel  string   `json:"channel,omitempty"`
		SocketID string   `json:"socket_id,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// The event data should not be larger than 10KB.
	if len(input.Data) > MAX_DATA_EVENT_SIZE {
		http.Error(w, "Request too large.", http.StatusRequestEntityTooLarge)
		return
	}

	// Trigger events

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// Allows fetching a hash of occupied channels (optionally filtered by prefix),
// and optionally one or more attributes for each channel.
//
// Notes:
// 'user_count' is the only attribute documented on the Pusher API
//
// Example:
// {
//   "channels": {
//     "presence-foobar": {
//       user_count: 42
//     },
//     "presence-another": {
//       user_count: 123
//     }
//   }
// }
//
// GET /apps/{app_id}/channels
func GetChannels(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	vars := mux.Vars(r)

	appID := vars["app_id"]
	filter := params.Get("filter_by_prefix")
	info := params.Get("info")

	attributes := strings.Split(info, ",")

	requestedUserCount := false

	for _, a := range attributes {
		if a == "user_count" {
			requestedUserCount = true
		}
	}

	// If an attribute such as user_count is requested, and the request is not limited
	// to presence channels, the API will return an error (400 code)
	if requestedUserCount && filter != "presence-" {
		http.Error(w, "Attribute user_count is restricted to presence channels", http.StatusBadRequest)
		return
	}

	app, err := Conf.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	channels := make(map[string]interface{})

	switch filter {
	case "presence-":
		for _, c := range app.PresenceChannels {
			if requestedUserCount {
				channels[c.ChannelID] = struct {
					UserCount int `json:"user_count"`
				}{
					c.totalUsers(),
				}
			} else {
				channels[c.ChannelID] = struct{}{}
			}
		}
	case "public-":
		for _, c := range app.PublicChannels {
			channels[c.ChannelID] = struct{}{}
		}
	case "private-":
		for _, c := range app.PrivateChannels {
			channels[c.ChannelID] = struct{}{}
		}
	default:
		for _, c := range app.AllChannels() {
			channels[c.ChannelID] = struct{}{}
		}
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}

// Fetch info for one channel
//
// Example:
// {
//   occupied: true,
//   user_count: 42,
//   subscription_count: 42
// }
//
// GET /apps/{app_id}/channels/{channel_name}
func GetChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	params := r.URL.Query()
	vars := mux.Vars(r)

	appID := vars["app_id"]
	app, err := Conf.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	channelName := vars["channel_name"]

	// Channel name could not be empty
	if strings.TrimSpace(channelName) == "" {
		http.Error(w, "Empty channel name", http.StatusBadRequest)
		return
	}

	info := params.Get("info")
	attributes := strings.Split(info, ",")

	// Attributes requested
	requestedUserCount := false
	requestedSubscriptionCount := false

	for _, a := range attributes {
		switch a {
		case "subscription_count":
			requestedSubscriptionCount = true
		case "user_count":
			requestedUserCount = true
		}
	}

	// Check the kind of channel

	channel, err := app.FindChannelByChannelID(channelName)

	// Channel exists?
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find a channel with id %s", channelName), http.StatusBadRequest)
		return
	}

	// If an attribute such as user_count is requested, and the request is not limited
	// to presence channels, the API will return an error (400 code)
	if requestedUserCount && !channel.isPresence() {
		http.Error(w, "Attribute user_count is restricted to presence channels", http.StatusBadRequest)
		return
	}

	// Output
	dtoChannel := struct {
		Occupied          bool `json:"occupied"`
		UserCount         int  `json:"user_count,omitempty"`
		SubscriptionCount int  `json:"subscription_count,omitempty"`
	}{Occupied: channel.isOccupied()}

	switch {
	case requestedSubscriptionCount && requestedUserCount:
		dtoChannel.UserCount = channel.totalUsers()
		dtoChannel.SubscriptionCount = channel.totalConnections()

	case requestedUserCount:
		dtoChannel.UserCount = channel.totalUsers()

	case requestedSubscriptionCount:
		dtoChannel.SubscriptionCount = channel.totalConnections()
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(dtoChannel); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}

// Allowed only to presence-channels
//
// Example:
// {
//  "users": [
//    { "id": 1 },
//    { "id": 2 }
//  ]
// }
//
// GET /apps/{app_id}/channels/{channel_name}/users
func GetChannelUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	appID := vars["app_id"]
	channelName := vars["channel_name"]

	isPresence := strings.HasPrefix(channelName, "presence-")

	if !isPresence {
		http.Error(w, "This api endpoint is restricted to presence channels.", http.StatusBadRequest)
		return
	}

	app, err := Conf.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	// Get the channel
	channel, err := app.FindChannelByChannelID(channelName)

	// Channel exists?
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find a channel with id %s", channelName), http.StatusBadRequest)
		return
	}

	result := make(map[string][]interface{})

	var users []interface{}

	for _, s := range channel.Connections {
		users = append(users, struct {
			Id int `json:"id"`
		}{s.Id})
	}

	result["users"] = users

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}
