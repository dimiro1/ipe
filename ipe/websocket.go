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

	goji "goji.io"

	"goji.io/pat"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"

	"github.com/dimiro1/ipe/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func handleMessages(conn *websocket.Conn, sessionID string, app *app) {
	var event struct {
		Event string `json:"event"`
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			handleError(conn, sessionID, app, err)
			return
		}

		if err := json.Unmarshal(message, &event); err != nil {
			emitWSError(newGenericReconnectImmediatelyError(), conn)
			return
		}

		log.Infof("websockets: Handling %s event", event.Event)

		switch event.Event {
		case "pusher:ping":
			onPing(conn)
		case "pusher:subscribe":
			onSubscribe(conn, sessionID, app, message)
		case "pusher:unsubscribe":
			onUnsubscribe(conn, sessionID, app, message)
		default:
			if utils.IsClientEvent(event.Event) {
				onClientEvent(conn, sessionID, app, message)
			}
		}
	}
}

func handleError(conn *websocket.Conn, sessionID string, app *app, err error) {
	log.Errorf("%+v", err)
	if err == io.EOF {
		onClose(sessionID, app)
	} else if _, ok := err.(*websocket.CloseError); ok {
		onClose(sessionID, app)
	} else {
		emitWSError(newGenericReconnectImmediatelyError(), conn)
	}
}

func onOpen(conn *websocket.Conn, r *http.Request, sessionID string, app *app) error {
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

func onClose(sessionID string, app *app) {
	app.Disconnect(sessionID)
}

func onPing(conn *websocket.Conn) {
	if err := conn.WriteJSON(newPongEvent()); err != nil {
		emitWSError(newGenericReconnectImmediatelyError(), conn)
	}
}

func onClientEvent(conn *websocket.Conn, sessionID string, app *app, message []byte) {
	if !app.UserEvents {
		emitWSError(newGenericError("To send client events, you must enable this feature in the Settings."), conn)
	}

	clientEvent := rawEvent{}

	if err := json.Unmarshal(message, &clientEvent); err != nil {
		log.Error(err)
		emitWSError(newGenericReconnectImmediatelyError(), conn)
		return
	}

	channel, err := app.FindChannelByChannelID(clientEvent.Channel)

	if err != nil {
		emitWSError(newGenericError(fmt.Sprintf("Could not find a channel with the id %s", clientEvent.Channel)), conn)
	}

	if !channel.IsPresenceOrPrivate() {
		emitWSError(newGenericError("Client event rejected - only supported on private and presence channels"), conn)
		return
	}

	if err := app.Publish(channel, clientEvent, sessionID); err != nil {
		log.Error(err)
		emitWSError(newGenericReconnectImmediatelyError(), conn)
		return
	}
}

func onUnsubscribe(conn *websocket.Conn, sessionID string, app *app, message []byte) {
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
		return
	}
}

func onSubscribe(conn *websocket.Conn, sessionID string, app *app, message []byte) {
	subscribeEvent := subscribeEvent{}

	if err := json.Unmarshal(message, &subscribeEvent); err != nil {
		emitWSError(newGenericReconnectImmediatelyError(), conn)
		return
	}

	connection, err := app.FindConnection(sessionID)

	if err != nil {
		emitWSError(newGenericReconnectImmediatelyError(), conn)
		return
	}

	channelName := strings.TrimSpace(subscribeEvent.Data.Channel)

	if !utils.IsChannelNameValid(channelName) {
		emitWSError(newGenericError("This channel name is not valid"), conn)
		return
	}

	isPresence := utils.IsPresenceChannel(channelName)
	isPrivate := utils.IsPrivateChannel(channelName)

	if isPresence || isPrivate {
		toSign := []string{connection.SocketID, channelName}

		if isPresence || len(subscribeEvent.Data.ChannelData) > 0 {
			toSign = append(toSign, subscribeEvent.Data.ChannelData)
		}

		if !validateAuthKey(subscribeEvent.Data.Auth, toSign, app) {
			emitWSError(newGenericError(fmt.Sprintf("Auth value for subscription to %s is invalid", channelName)), conn)
			return
		}
	}

	channel := app.FindOrCreateChannelByChannelID(channelName)
	log.Info(subscribeEvent.Data.ChannelData)

	if err := app.Subscribe(channel, connection, subscribeEvent.Data.ChannelData); err != nil {
		emitWSError(newGenericReconnectImmediatelyError(), conn)
	}
}

func validateAuthKey(givenAuthKey string, toSign []string, app *app) bool {
	expectedAuthKey := fmt.Sprintf("%s:%s", app.Key, utils.HashMAC([]byte(strings.Join(toSign, ":")), []byte(app.Secret)))
	return givenAuthKey == expectedAuthKey
}

// Emit an Websocket ErrorEvent
func emitWSError(err error, conn *websocket.Conn) {
	e, ok := err.(websocketError)

	if !ok {
		log.Error(err)
		return
	}

	event := newErrorEvent(e.GetCode(), e.GetMsg())

	if err := conn.WriteJSON(event); err != nil {
		log.Error(err)
	}
}

func newWebsocketHandler(DB db) goji.Handler {
	return &websocketHandler{DB}
}

type websocketHandler struct {
	DB db
}

// Websocket GET /app/{key}
func (h *websocketHandler) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	appKey := pat.Param(ctx, "key")

	app, err := h.DB.GetAppByKey(appKey)

	if err != nil {
		log.Error(err)
		emitWSError(newApplicationDoesNotExistsError(), conn)
		return
	}

	sessionID := utils.GenerateSessionID()

	if err := onOpen(conn, r, sessionID, app); err != nil {
		emitWSError(err, conn)
		return
	}

	handleMessages(conn, sessionID, app)
}
