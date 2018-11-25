// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package websockets

type websocketError struct {
	Code int
	Msg  string
}

var (
	applicationOnlyAcceptsSSL  = &websocketError{Code: 4000, Msg: "Application only accepts SSL connections, reconnect using wss://"}
	applicationDoesNotExists   = &websocketError{Code: 4001, Msg: "Could not found an app with the given key"}
	applicationDisabled        = &websocketError{Code: 4003, Msg: "Application disabled"}
	invalidVersionStringFormat = &websocketError{Code: 4006, Msg: "Invalid version string format"}
	unsupportedProtocolVersion = &websocketError{Code: 4007, Msg: "Unsupported protocol version"}
	noProtocolVersionSupplied  = &websocketError{Code: 4008, Msg: "No protocol version supplied"}
	reconnectImmediately       = &websocketError{Code: 4200, Msg: "Generic reconnect immediately"}
)
