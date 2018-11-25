// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package websockets

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

	"ipe/app"
	"ipe/connection"
	"ipe/events"
	"ipe/storage"
	"ipe/utils"
)

// Only this version is supported
const supportedProtocolVersion = 7

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type Websocket struct {
	storage storage.Storage
}

func NewWebsocket(storage storage.Storage) *Websocket {
	return &Websocket{storage: storage}
}

// Websocket GET /app/{key}
func (h *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Errorf("closing the websocket connection %+v", err)
			}
		}
	}()

	if err != nil {
		log.Error(err)
		return
	}

	var (
		pathVars = mux.Vars(r)
		appKey   = pathVars["key"]
	)

	_app, err := h.storage.GetAppByKey(appKey)

	if err != nil {
		log.Error(err)
		emitError(applicationDoesNotExists, conn)
		return
	}

	sessionID := utils.GenerateSessionID()

	if err := onOpen(conn, r, sessionID, _app); err != nil {
		emitError(err, conn)
		return
	}

	handleMessages(conn, sessionID, _app)
}

func handleMessages(conn *websocket.Conn, sessionID string, app *app.Application) {
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
			emitError(reconnectImmediately, conn)
			return
		}

		log.Infof("websocket: Handling %s event", event.Event)

		switch event.Event {
		case "pusher:ping":
			handlePing(conn)
		case "pusher:subscribe":
			handleSubscribe(conn, sessionID, app, message)
		case "pusher:unsubscribe":
			handleUnsubscribe(conn, sessionID, app, message)
		default:
			if utils.IsClientEvent(event.Event) {
				handleClientEvent(conn, sessionID, app, message)
			}
		}
	}
}

// Emit an Websocket ErrorEvent
func emitError(err *websocketError, conn *websocket.Conn) {
	event := events.NewError(err.Code, err.Msg)

	if err := conn.WriteJSON(event); err != nil {
		log.Error(err)
	}
}

func handleError(conn *websocket.Conn, sessionID string, app *app.Application, err error) {
	log.Errorf("%+v", err)
	if err == io.EOF {
		onClose(sessionID, app)
	} else if _, ok := err.(*websocket.CloseError); ok {
		onClose(sessionID, app)
	} else {
		emitError(reconnectImmediately, conn)
	}
}

func onOpen(conn *websocket.Conn, r *http.Request, sessionID string, app *app.Application) *websocketError {
	var (
		queryVars   = r.URL.Query()
		strProtocol = queryVars.Get("protocol")
	)

	protocol, err := strconv.Atoi(strProtocol)
	if err != nil {
		return invalidVersionStringFormat
	}

	switch {
	case strings.TrimSpace(strProtocol) == "":
		return noProtocolVersionSupplied
	case protocol != supportedProtocolVersion:
		return unsupportedProtocolVersion
	case app.ApplicationDisabled:
		return applicationDisabled
	case app.OnlySSL:
		if r.TLS == nil {
			return applicationOnlyAcceptsSSL
		}
	}

	// Create the new Subscriber
	_connection := connection.New(sessionID, conn)
	app.Connect(_connection)

	// Everything went fine.
	if err := conn.WriteJSON(events.NewConnectionEstablished(_connection.SocketID)); err != nil {
		return reconnectImmediately
	}

	return nil
}

func onClose(sessionID string, app *app.Application) {
	app.Disconnect(sessionID)
}

func handlePing(conn *websocket.Conn) {
	if err := conn.WriteJSON(events.NewPong()); err != nil {
		emitError(reconnectImmediately, conn)
	}
}

func handleClientEvent(conn *websocket.Conn, sessionID string, app *app.Application, message []byte) {
	if !app.UserEvents {
		emitError(&websocketError{Code: 0, Msg: "To send client events, you must enable this feature in the Settings."}, conn)
	}

	clientEvent := events.Raw{}

	if err := json.Unmarshal(message, &clientEvent); err != nil {
		log.Error(err)
		emitError(reconnectImmediately, conn)
		return
	}

	channel, err := app.FindChannelByChannelID(clientEvent.Channel)

	if err != nil {
		emitError(&websocketError{Code: 0, Msg: fmt.Sprintf("Could not find a channel with the id %s", clientEvent.Channel)}, conn)
		return
	}

	if !channel.IsPresenceOrPrivate() {
		emitError(&websocketError{Code: 0, Msg: "Client event rejected - only supported on private and presence channels"}, conn)
		return
	}

	if err := app.Publish(channel, clientEvent, sessionID); err != nil {
		log.Error(err)
		emitError(reconnectImmediately, conn)
		return
	}
}

func handleUnsubscribe(conn *websocket.Conn, sessionID string, app *app.Application, message []byte) {
	unsubscribeEvent := events.Unsubscribe{}

	if err := json.Unmarshal(message, &unsubscribeEvent); err != nil {
		emitError(reconnectImmediately, conn)
		return
	}

	_connection, err := app.FindConnection(sessionID)

	if err != nil {
		emitError(&websocketError{Code: 0, Msg: fmt.Sprintf("Could not find a connection with the id %s", sessionID)}, conn)
		return
	}

	channel, err := app.FindChannelByChannelID(unsubscribeEvent.Data.Channel)

	if err != nil {
		emitError(&websocketError{Code: 0, Msg: fmt.Sprintf("Could not find a channel with the id %s", unsubscribeEvent.Data.Channel)}, conn)
		return
	}

	if err := app.Unsubscribe(channel, _connection); err != nil {
		emitError(reconnectImmediately, conn)
		return
	}
}

func handleSubscribe(conn *websocket.Conn, sessionID string, app *app.Application, message []byte) {
	subscribeEvent := events.Subscribe{}

	if err := json.Unmarshal(message, &subscribeEvent); err != nil {
		emitError(reconnectImmediately, conn)
		return
	}

	_connection, err := app.FindConnection(sessionID)

	if err != nil {
		emitError(reconnectImmediately, conn)
		return
	}

	channelName := strings.TrimSpace(subscribeEvent.Data.Channel)

	if !utils.IsChannelNameValid(channelName) {
		emitError(&websocketError{Code: 0, Msg: "This channel name is not valid"}, conn)
		return
	}

	isPresence := utils.IsPresenceChannel(channelName)
	isPrivate := utils.IsPrivateChannel(channelName)

	if isPresence || isPrivate {
		toSign := []string{_connection.SocketID, channelName}

		if isPresence || len(subscribeEvent.Data.ChannelData) > 0 {
			toSign = append(toSign, subscribeEvent.Data.ChannelData)
		}

		if !validateAuthKey(subscribeEvent.Data.Auth, toSign, app) {
			emitError(&websocketError{Code: 0, Msg: fmt.Sprintf("Auth value for subscription to %s is invalid", channelName)}, conn)
			return
		}
	}

	channel := app.FindOrCreateChannelByChannelID(channelName)
	log.Info(subscribeEvent.Data.ChannelData)

	if err := app.Subscribe(channel, _connection, subscribeEvent.Data.ChannelData); err != nil {
		emitError(reconnectImmediately, conn)
	}
}

func validateAuthKey(givenAuthKey string, toSign []string, app *app.Application) bool {
	expectedAuthKey := fmt.Sprintf("%s:%s", app.Key, utils.HashMAC([]byte(strings.Join(toSign, ":")), []byte(app.Secret)))
	return givenAuthKey == expectedAuthKey
}
