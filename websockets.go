// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handle open connection.
func onOpen(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, session *sessions.Session, app *App) WebsocketError {
	params := r.URL.Query()
	p := params.Get("protocol")

	protocol, err := strconv.Atoi(p)

	if err != nil {
		return NewInvalidVersionStringFormatError()
	}

	switch {
	case strings.TrimSpace(p) == "":
		return NewNoProtocolVersionSuppliedError()
	case protocol != SUPPORTED_PROTOCOL_VERSION:
		return NewUnsupportedProtocolVersionError()
	case app.ApplicationDisabled:
		return NewApplicationDisabledError()
	case r.TLS != nil:
		if app.OnlySSL {
			return NewApplicationOnlyAccepsSSLError()
		}
	}

	// Create the new connection
	connection := NewConnection(session.ID, "", conn)
	app.AddConnection(connection)

	// Everything went fine. Huhu.
	if err := conn.WriteJSON(NewConnectionEstablishedEvent(connection.SocketID)); err != nil {
		return NewGenericReconnectImmediatelyError()
	}

	return nil
}

// Handle messages
func onMessage(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, session *sessions.Session, app *App) WebsocketError {
	var event struct {
		Event string `json:"event"`
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			return NewGenericReconnectImmediatelyError()
		}

		if err := json.Unmarshal(message, &event); err != nil {
			return NewGenericReconnectImmediatelyError()
		}

		switch event.Event {
		case "pusher:ping":
			if err := conn.WriteJSON(NewPongEvent()); err != nil {
				return NewGenericReconnectImmediatelyError()
			}
		case "pusher:subscribe":
			subscribeEvent := SubscribeEvent{}

			if err := json.Unmarshal(message, &subscribeEvent); err != nil {
				return NewGenericReconnectImmediatelyError()
			}

			connection, err := app.FindConnection(session.ID)

			if err != nil {
				return NewGenericReconnectImmediatelyError()
			}

			channelName := strings.TrimSpace(subscribeEvent.Data.Channel)

			// Authentication
			if strings.HasPrefix(channelName, "presence-") {
				toSign := fmt.Sprintf("%s:%s:%s", connection.SocketID, channelName, subscribeEvent.Data.ChannelData)

				if subscribeEvent.Data.Auth != HashMAC([]byte(toSign), []byte(app.Secret)) {
					return NewGenericError(fmt.Sprintf("Auth value for subscription to %s is invalid", channelName))
				}
			} else if strings.HasPrefix(channelName, "private-") {
				toSign := fmt.Sprintf("%s:%s", connection.SocketID, channelName)

				if subscribeEvent.Data.Auth != HashMAC([]byte(toSign), []byte(app.Secret)) {
					return NewGenericError(fmt.Sprintf("Auth value for subscription to %s is invalid", channelName))
				}
			}

			channel := app.FindOrCreateChannelByChannelID(channelName, subscribeEvent.Data.ChannelData)
			channel.Subscribe(connection)

			if err := conn.WriteJSON(NewSubscriptionSucceededEvent(channel.ChannelID)); err != nil {
				return NewGenericReconnectImmediatelyError()
			}
		case "pusher:unsubscribe":
			unsubscribeEvent := UnsubscribeEvent{}

			if err := json.Unmarshal(message, &unsubscribeEvent); err != nil {
				return NewGenericReconnectImmediatelyError()
			}

			connection, err := app.FindConnection(session.ID)

			if err != nil {
				return NewGenericError(fmt.Sprintf("Could not find a connection with the id %s", session.ID))
			}

			channel, err := app.FindChannelByChannelID(unsubscribeEvent.Data.Channel)

			if err != nil {
				return NewGenericError(fmt.Sprintf("Could not find a channel with the id %s", unsubscribeEvent.Data.Channel))
			}

			if err := channel.Unsubscribe(connection); err != nil {
				return NewGenericReconnectImmediatelyError()
			}
		}

		return nil
		// Client Events
	}
}

// Websocket GET /app/{key}
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		emitWSError(NewGenericReconnectImmediatelyError(), conn)
		return
	}

	var store = sessions.NewFilesystemStore("", []byte(Conf.SessionSecret))
	session, err := store.Get(r, Conf.SessionName)

	if err != nil {
		log.Println(err)
		emitWSError(NewGenericReconnectImmediatelyError(), conn)
		return
	}

	if err := session.Save(r, w); err != nil {
		log.Println(err)
		emitWSError(NewGenericReconnectImmediatelyError(), conn)
		return
	}

	vars := mux.Vars(r)
	appKey := vars["key"]

	app, err := Conf.GetAppByKey(appKey)

	if err != nil {
		log.Println(err)
		emitWSError(NewApplicationDoesNotExistsError(), conn)
		return
	}

	if err := onOpen(conn, w, r, session, app); err != nil {
		emitWSError(err, conn)
		return
	}

	if err := onMessage(conn, w, r, session, app); err != nil {
		emitWSError(err, conn)

		// Find the connection in app and destroy it
		return
	}
}

// Emit an Websocket ErrorEvent
func emitWSError(err WebsocketError, conn *websocket.Conn) {
	event := NewErrorEvent(err.GetCode(), err.GetMsg())

	if err := conn.WriteJSON(event); err != nil {
		log.Println(err)
	}

	conn.Close()
}
