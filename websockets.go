// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/dimiro1/ipe/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handle open Subscriber.
func onOpen(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, sessionID string, app *App) WebsocketError {
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

	// Create the new Subscriber
	connection := NewConnection(sessionID, conn)
	app.Connect(connection)

	// Everything went fine. Huhu.
	if err := conn.WriteJSON(NewConnectionEstablishedEvent(connection.SocketID)); err != nil {
		return NewGenericReconnectImmediatelyError()
	}

	return nil
}

// Handle the close event
func onClose(sessionID string, app *App) {
	app.Disconnect(sessionID)
}

// Handle messages
//
// If there is an unrecoverable error then break the loop,
// otherwise just keep going.
func onMessage(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, sessionID string, app *App) {
	var event struct {
		Event string `json:"event"`
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Errorf("%+v", err)
			switch err {
			case io.EOF:
				onClose(sessionID, app)
			default:
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
			}
			break
		}

		if err := json.Unmarshal(message, &event); err != nil {
			emitWSError(NewGenericReconnectImmediatelyError(), conn)
			break
		}

		log.Infof("websockets: Handling %s event", event.Event)

		switch event.Event {
		case "pusher:ping":
			if err := conn.WriteJSON(NewPongEvent()); err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
			}
		case "pusher:subscribe":
			subscribeEvent := SubscribeEvent{}

			if err := json.Unmarshal(message, &subscribeEvent); err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
				break
			}

			connection, err := app.FindConnection(sessionID)

			if err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
				break
			}

			channelName := strings.TrimSpace(subscribeEvent.Data.Channel)

			isPresence := strings.HasPrefix(channelName, "presence-")
			isPrivate := strings.HasPrefix(channelName, "private-")

			if isPresence || isPrivate {
				toSign := []string{connection.SocketID, channelName}

				if isPresence {
					toSign = append(toSign, subscribeEvent.Data.ChannelData)
				}

				expectedAuthKey := fmt.Sprintf("%s:%s", app.Key, utils.HashMAC([]byte(strings.Join(toSign, ":")), []byte(app.Secret)))
				if subscribeEvent.Data.Auth != expectedAuthKey {
					emitWSError(NewGenericError(fmt.Sprintf("Auth value for subscription to %s is invalid", channelName)), conn)
					continue
				}
			}

			channel := app.FindOrCreateChannelByChannelID(channelName)
			log.Info(subscribeEvent.Data.ChannelData)

			if err := app.Subscribe(channel, connection, subscribeEvent.Data.ChannelData); err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
			}
		case "pusher:unsubscribe":
			unsubscribeEvent := UnsubscribeEvent{}

			if err := json.Unmarshal(message, &unsubscribeEvent); err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
			}

			connection, err := app.FindConnection(sessionID)

			if err != nil {
				emitWSError(NewGenericError(fmt.Sprintf("Could not find a connection with the id %s", sessionID)), conn)
			}

			channel, err := app.FindChannelByChannelID(unsubscribeEvent.Data.Channel)

			if err != nil {
				emitWSError(NewGenericError(fmt.Sprintf("Could not find a channel with the id %s", unsubscribeEvent.Data.Channel)), conn)
			}

			if err := app.Unsubscribe(channel, connection); err != nil {
				emitWSError(NewGenericReconnectImmediatelyError(), conn)
				break
			}
		default: // CLient Events ??
			// see http://pusher.com/docs/client_api_guide/client_events#trigger-events
			if strings.HasPrefix(event.Event, "client-") {
				if !app.UserEvents {
					emitWSError(NewGenericError("To send client events, you must enable this feature in the Settings."), conn)
				}

				clientEvent := RawEvent{}

				if err := json.Unmarshal(message, &clientEvent); err != nil {
					log.Error(err)
					emitWSError(NewGenericReconnectImmediatelyError(), conn)
					break
				}

				channel, err := app.FindChannelByChannelID(clientEvent.Channel)

				if !channel.IsPresenceOrPrivate() {
					emitWSError(NewGenericError("Client event rejected - only supported on private and presence channels"), conn)
					break
				}

				if err != nil {
					emitWSError(NewGenericError(fmt.Sprintf("Could not find a channel with the id %s", clientEvent.Channel)), conn)
				}

				if err := app.Publish(channel, clientEvent, sessionID); err != nil {
					log.Error(err)
					emitWSError(NewGenericReconnectImmediatelyError(), conn)
					break
				}
			}

		} // switch
	} // For
}

// Websocket GET /app/{key}
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	if err != nil {
		log.Error(err)
		emitWSError(NewGenericReconnectImmediatelyError(), conn)
		return
	}

	vars := mux.Vars(r)
	appKey := vars["key"]

	app, err := Conf.GetAppByKey(appKey)

	if err != nil {
		log.Error(err)
		emitWSError(NewApplicationDoesNotExistsError(), conn)
		return
	}

	sessionID := utils.RandomHash()

	if err := onOpen(conn, w, r, sessionID, app); err != nil {
		emitWSError(err, conn)
		return
	}

	onMessage(conn, w, r, sessionID, app)
}

// Emit an Websocket ErrorEvent
func emitWSError(err WebsocketError, conn *websocket.Conn) {

	event := NewErrorEvent(err.GetCode(), err.GetMsg())

	if err := conn.WriteJSON(event); err != nil {
		log.Error(err)
	}
}
