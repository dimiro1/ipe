// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"

	"github.com/dimiro1/ipe/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handle open Subscriber.
func onOpen(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, sessionID string, app *app) websocketError {
	params := r.URL.Query()
	p := params.Get("protocol")

	protocol, err := strconv.Atoi(p)

	if err != nil {
		return newInvalidVersionStringFormatError()
	}

	switch {
	case strings.TrimSpace(p) == "":
		return newNoProtocolVersionSuppliedError()
	case protocol != supportedProtocolVersion:
		return newUnsupportedProtocolVersionError()
	case app.ApplicationDisabled:
		return newApplicationDisabledError()
	case app.OnlySSL:
		if r.TLS == nil {
			return newApplicationOnlyAccepsSSLError()
		}
	}

	// Create the new Subscriber
	connection := newConnection(sessionID, conn)
	app.Connect(connection)

	// Everything went fine. Huhu.
	if err := conn.WriteJSON(newConnectionEstablishedEvent(connection.SocketID)); err != nil {
		return newGenericReconnectImmediatelyError()
	}

	return nil
}

// Handle the close event
func onClose(sessionID string, app *app) {
	app.Disconnect(sessionID)
}

// Handle messages
//
// If there is an unrecoverable error then break the loop,
// otherwise just keep going.
func onMessage(conn *websocket.Conn, w http.ResponseWriter, r *http.Request, sessionID string, app *app) {
	var event struct {
		Event string `json:"event"`
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Errorf("%+v", err)
			if _, ok := err.(*websocket.CloseError); ok {
				onClose(sessionID, app)
			} else {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
			}
			break
		}

		if err := json.Unmarshal(message, &event); err != nil {
			emitWSError(newGenericReconnectImmediatelyError(), conn)
			break
		}

		log.Infof("websockets: Handling %s event", event.Event)

		switch event.Event {
		case "pusher:ping":
			if err := conn.WriteJSON(newPongEvent()); err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
			}
		case "pusher:subscribe":
			subscribeEvent := subscribeEvent{}

			if err := json.Unmarshal(message, &subscribeEvent); err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
				break
			}

			connection, err := app.FindConnection(sessionID)

			if err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
				break
			}

			channelName := strings.TrimSpace(subscribeEvent.Data.Channel)

			if !utils.IsChannelNameValid(channelName) {
				emitWSError(newGenericError(fmt.Sprintf("This channel name is not valid")), conn)
				break
			}

			isPresence := utils.IsPresenceChannel(channelName)
			isPrivate := utils.IsPrivateChannel(channelName)

			if isPresence || isPrivate {
				toSign := []string{connection.SocketID, channelName}

				if isPresence || len(subscribeEvent.Data.ChannelData) > 0 {
					toSign = append(toSign, subscribeEvent.Data.ChannelData)
				}

				expectedAuthKey := fmt.Sprintf("%s:%s", app.Key, utils.HashMAC([]byte(strings.Join(toSign, ":")), []byte(app.Secret)))
				if subscribeEvent.Data.Auth != expectedAuthKey {
					emitWSError(newGenericError(fmt.Sprintf("Auth value for subscription to %s is invalid", channelName)), conn)
					continue
				}
			}

			channel := app.FindOrCreateChannelByChannelID(channelName)
			log.Info(subscribeEvent.Data.ChannelData)

			if err := app.Subscribe(channel, connection, subscribeEvent.Data.ChannelData); err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
			}
		case "pusher:unsubscribe":
			unsubscribeEvent := unsubscribeEvent{}

			if err := json.Unmarshal(message, &unsubscribeEvent); err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
			}

			connection, err := app.FindConnection(sessionID)

			if err != nil {
				emitWSError(newGenericError(fmt.Sprintf("Could not find a connection with the id %s", sessionID)), conn)
			}

			channel, err := app.FindChannelByChannelID(unsubscribeEvent.Data.Channel)

			if err != nil {
				emitWSError(newGenericError(fmt.Sprintf("Could not find a channel with the id %s", unsubscribeEvent.Data.Channel)), conn)
			}

			if err := app.Unsubscribe(channel, connection); err != nil {
				emitWSError(newGenericReconnectImmediatelyError(), conn)
				break
			}
		default: // CLient Events ??
			// see http://pusher.com/docs/client_api_guide/client_events#trigger-events
			if strings.HasPrefix(event.Event, "client-") {
				if !app.UserEvents {
					emitWSError(newGenericError("To send client events, you must enable this feature in the Settings."), conn)
				}

				clientEvent := rawEvent{}

				if err := json.Unmarshal(message, &clientEvent); err != nil {
					log.Error(err)
					emitWSError(newGenericReconnectImmediatelyError(), conn)
					break
				}

				channel, err := app.FindChannelByChannelID(clientEvent.Channel)

				if err != nil {
					emitWSError(newGenericError(fmt.Sprintf("Could not find a channel with the id %s", clientEvent.Channel)), conn)
				}

				if !channel.IsPresenceOrPrivate() {
					emitWSError(newGenericError("Client event rejected - only supported on private and presence channels"), conn)
					break
				}

				if err := app.Publish(channel, clientEvent, sessionID); err != nil {
					log.Error(err)
					emitWSError(newGenericReconnectImmediatelyError(), conn)
					break
				}
			}

		} // switch
	} // For
}

// Websocket GET /app/{key}
func wsHandler(ctx *applicationContext, p params, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		log.Error(err)
		return
	}

	appKey := p.Get("key")

	app, err := ctx.DB.GetAppByKey(appKey)

	if err != nil {
		log.Error(err)
		emitWSError(newApplicationDoesNotExistsError(), conn)
		return
	}

	sessionID := utils.GenerateSessionID()

	if err := onOpen(conn, w, r, sessionID, app); err != nil {
		emitWSError(err, conn)
		return
	}

	onMessage(conn, w, r, sessionID, app)
}

// Emit an Websocket ErrorEvent
func emitWSError(err websocketError, conn *websocket.Conn) {

	event := newErrorEvent(err.GetCode(), err.GetMsg())

	if err := conn.WriteJSON(event); err != nil {
		log.Error(err)
	}
}
