// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "fmt"

// Base struct
type websocketError struct {
	Code *int
	Msg  string
}

func (e websocketError) GetCode() *int {
	return e.Code
}

func (e websocketError) GetMsg() string {
	return e.Msg
}

func (e websocketError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func newWebsocketError(code int, msg string) websocketError {
	return websocketError{Code: &code, Msg: msg}
}

var (
	// Unsupprted protocol version
	unsupportedProtocolVersionError = newWebsocketError(4007, "Unsupported protocol version")

	// The application does not exists
	// See the configuration file
	applicationDoesNotExistsError = newWebsocketError(4001, "Could not found an app with the given key")

	// The user did not send the protocol version
	noProtocolVersionSuppliedError = newWebsocketError(4008, "No protocol version supplied")

	// When the application is disabled.
	// See the configuration file
	applicationDisabledError = newWebsocketError(4003, "Application disabled")

	// When the application only accepts SSL connections
	applicationOnlyAccepsSSLError = newWebsocketError(4000, "Application only accepts SSL connections, reconnect using wss://")

	// When the user send an invalid version
	invalidVersionStringFormatError = newWebsocketError(4006, "Invalid version string format")

	// Used when the error was internal
	// * Decoding json
	// * Writing to output
	genericReconnectImmediatelyError = newWebsocketError(4200, "Generic reconnect immediately")

	// When pusher wants to send an Generic error, it only send the message, the code become nil
	// Currently I do not know how to send nil, so I send GENERIC_ERROR
	genericError = newWebsocketError(0, "Generic Error")

	disabledClientEventsError = websocketError{Msg: "To send client events, you must enable this feature in the Settings."}

	couldNotFoundChannelError = websocketError{Msg: "Could not find a channel with the given id"}
)
