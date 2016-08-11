// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	goji "goji.io"

	"goji.io/pat"

	log "github.com/golang/glog"
	"golang.org/x/net/context"

	"github.com/dimiro1/ipe/utils"
)

// Prepare Querystring
func prepareQueryString(params url.Values) string {
	var keys []string

	for key := range params {
		keys = append(keys, strings.ToLower(key))
	}

	sort.Strings(keys)

	var pieces []string

	for _, key := range keys {
		pieces = append(pieces, key+"="+params.Get(key))
	}

	return strings.Join(pieces, "&")
}

// Authenticate pusher
// see: https://gist.github.com/mloughran/376898
//
// The signature is a HMAC SHA256 hex digest.
// This is generated by signing a string made up of the following components concatenated with newline characters \n.
//
//  * The uppercase request method (e.g. POST)
//  * The request path (e.g. /some/resource)
//  * The query parameters sorted by key, with keys converted to lowercase, then joined as in the query string.
//    Note that the string must not be url escaped (e.g. given the keys auth_key: foo, Name: Something else, you get auth_key=foo&name=Something else)
func restAuthenticationHandler(DB db, h goji.Handler) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		appID := pat.Param(ctx, "app_id")

		app, err := DB.GetAppByAppID(appID)

		if err != nil {
			log.Error(err)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		query := r.URL.Query()

		signature := query.Get("auth_signature")
		query.Del("auth_signature")

		queryString := prepareQueryString(query)

		toSign := strings.ToUpper(r.Method) + "\n" + r.URL.Path + "\n" + queryString

		if utils.HashMAC([]byte(toSign), []byte(app.Secret)) == signature {
			h.ServeHTTPC(ctx, w, r)
		} else {
			log.Error("Not authorized")
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		}
	}
}

// Check if the application is disabled
func restCheckAppDisabledHandler(DB db, h goji.Handler) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		appID := pat.Param(ctx, "app_id")

		currentApp, err := DB.GetAppByAppID(appID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusForbidden)
			return
		}

		if currentApp.ApplicationDisabled {
			http.Error(w, "Application disabled", http.StatusForbidden)
			return
		}

		h.ServeHTTPC(ctx, w, r)
	}
}

func recoverHandle(h goji.Handler) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Please verify the url parameters error was: %s", r)
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}
		}()
		h.ServeHTTPC(ctx, w, r)
	}
}

// commonHandlers combine restCheckAppDisabledHandler and restAuthenticationHandler handlers
func commonHandlers(DB db, h goji.Handler) goji.HandlerFunc {
	return recoverHandle(restCheckAppDisabledHandler(DB, restAuthenticationHandler(DB, h)))
}

type postEvents struct{ DB db }

// ServeHTTPC An event consists of a name and data (typically JSON) which may be sent to all subscribers to a particular channel or channels.
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
func (h *postEvents) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := pat.Param(ctx, "app_id")

	app, err := h.DB.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	var input struct {
		Name     string          `json:"name"`
		Data     json.RawMessage `json:"data"`
		Channels []string        `json:"channels,omitempty"`
		Channel  string          `json:"channel,omitempty"`
		SocketID string          `json:"socket_id,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// The event data should not be larger than 10KB.
	if len(input.Data) > maxDataEventSize {
		http.Error(w, "Request too large.", http.StatusRequestEntityTooLarge)
		return
	}

	log.Info(input.Channels)
	if len(input.Channel) > 0 && len(input.Channels) == 0 {
		input.Channels = append(input.Channels, input.Channel)
	}

	for _, c := range input.Channels {
		channel := app.FindOrCreateChannelByChannelID(c)

		app.Publish(channel, rawEvent{Event: input.Name, Channel: c, Data: input.Data}, input.SocketID)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

type getChannels struct{ DB db }

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
func (h *getChannels) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	appID := pat.Param(ctx, "app_id")
	filter := query.Get("filter_by_prefix")
	info := query.Get("info")

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

	app, err := h.DB.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	channels := make(map[string]interface{})

	switch filter {
	case "presence-":
		for _, c := range app.PresenceChannels() {
			if requestedUserCount {
				channels[c.ChannelID] = struct {
					UserCount int `json:"user_count"`
				}{
					c.TotalUsers(),
				}
			} else {
				channels[c.ChannelID] = struct{}{}
			}
		}
	case "public-":
		for _, c := range app.PublicChannels() {
			channels[c.ChannelID] = struct{}{}
		}
	case "private-":
		for _, c := range app.PrivateChannels() {
			channels[c.ChannelID] = struct{}{}
		}
	default:
		for _, c := range app.Channels {
			channels[c.ChannelID] = struct{}{}
		}
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	js := make(map[string]interface{}, 1)
	js["channels"] = channels

	if err := json.NewEncoder(w).Encode(js); err != nil {
		log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type getChannel struct{ DB db }

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
func (h *getChannel) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	query := r.URL.Query()

	appID := pat.Param(ctx, "app_id")
	app, err := h.DB.GetAppByAppID(appID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusBadRequest)
	}

	channelName := pat.Param(ctx, "channel_name")

	// Channel name could not be empty
	if strings.TrimSpace(channelName) == "" {
		http.Error(w, "Empty channel name", http.StatusBadRequest)
		return
	}

	info := query.Get("info")
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

	channel, err := app.FindChannelByChannelID(channelName)

	// Channel exists?
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find a channel with id %s", channelName), http.StatusBadRequest)
		return
	}

	// If an attribute such as user_count is requested, and the request is not limited
	// to presence channels, the API will return an error (400 code)
	if requestedUserCount && !channel.IsPresence() {
		http.Error(w, "Attribute user_count is restricted to presence channels", http.StatusBadRequest)
		return
	}

	// Output
	dtoChannel := struct {
		Occupied          bool `json:"occupied"`
		UserCount         int  `json:"user_count,omitempty"`
		SubscriptionCount int  `json:"subscription_count,omitempty"`
	}{Occupied: channel.IsOccupied()}

	switch {
	case requestedSubscriptionCount && requestedUserCount:
		dtoChannel.UserCount = channel.TotalUsers()
		dtoChannel.SubscriptionCount = channel.TotalSubscriptions()

	case requestedUserCount:
		dtoChannel.UserCount = channel.TotalUsers()

	case requestedSubscriptionCount:
		dtoChannel.SubscriptionCount = channel.TotalSubscriptions()
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(dtoChannel); err != nil {
		log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type getChannelUsers struct{ DB db }

// Allowed only for presence-channels
//
// Example:
// {
//  "users": [
//    { "id": "1" },
//    { "id": "2" }
//  ]
// }
//
// GET /apps/{app_id}/channels/{channel_name}/users
func (h *getChannelUsers) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := pat.Param(ctx, "app_id")
	channelName := pat.Param(ctx, "channel_name")

	isPresence := utils.IsPresenceChannel(channelName)

	if !isPresence {
		http.Error(w, "This api endpoint is restricted to presence channels.", http.StatusBadRequest)
		return
	}

	app, err := h.DB.GetAppByAppID(appID)

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

	for _, s := range channel.Subscriptions {
		users = append(users, struct {
			ID string `json:"id"`
		}{s.ID})
	}

	result["users"] = users

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Error(err)
	}
}
