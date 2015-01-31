// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

// Base interface
type WebsocketError interface {
	GetCode() int
	GetMsg() string
}

// Base struct
type BaseWebsocketError struct {
	Code int
	Msg  string
}

func (e BaseWebsocketError) GetCode() int {
	return e.Code
}

func (e BaseWebsocketError) GetMsg() string {
	return e.Msg
}

// Unsupprted protocol version
type UnsupportedProtocolVersionError struct {
	BaseWebsocketError
}

func newUnsupportedProtocolVersionError() UnsupportedProtocolVersionError {
	return UnsupportedProtocolVersionError{
		BaseWebsocketError{Code: UNSUPPORTED_PROTOCOL_VERSION, Msg: "Unsupported protocol version"},
	}
}

// The application does not exists
// See the configuration file
type ApplicationDoesNotExistsError struct {
	BaseWebsocketError
}

func newApplicationDoesNotExistsError() ApplicationDoesNotExistsError {
	return ApplicationDoesNotExistsError{
		BaseWebsocketError{Code: APPLICATION_DOES_NOT_EXISTS, Msg: "Could not found an app with the given key"},
	}
}

// The user did not send the protocol version
type NoProtocolVersionSuppliedError struct {
	BaseWebsocketError
}

func newNoProtocolVersionSuppliedError() NoProtocolVersionSuppliedError {
	return NoProtocolVersionSuppliedError{
		BaseWebsocketError{Code: NO_PROTOCOL_VERSION_SUPPLIED, Msg: "No protocol version supplied"},
	}
}

// When the application is disabled.
// See the configuration file
type ApplicationDisabledError struct {
	BaseWebsocketError
}

func newApplicationDisabledError() NoProtocolVersionSuppliedError {
	return NoProtocolVersionSuppliedError{
		BaseWebsocketError{Code: APPLICATION_DISABLED, Msg: "Application disabled"},
	}
}

// When the application only accepts SSL connections
type ApplicationOnlyAccepsSSLError struct {
	BaseWebsocketError
}

func newApplicationOnlyAccepsSSLError() ApplicationOnlyAccepsSSLError {
	return ApplicationOnlyAccepsSSLError{
		BaseWebsocketError{Code: APPLICATION_ONLY_ACCEPTS_SSL, Msg: "Application only accepts SSL connections, reconnect using wss://"},
	}
}

// When the user send an invalid version
type InvalidVersionStringFormatError struct {
	BaseWebsocketError
}

func newInvalidVersionStringFormatError() InvalidVersionStringFormatError {
	return InvalidVersionStringFormatError{
		BaseWebsocketError{Code: INVALID_VERSION_STRING_FORMAT, Msg: "Invalid version string format"},
	}
}

// Used when the error was internal
// * Decoding json
// * Writing to output
type GenericReconnectImmediatelyError struct {
	BaseWebsocketError
}

func newGenericReconnectImmediatelyError() GenericReconnectImmediatelyError {
	return GenericReconnectImmediatelyError{
		BaseWebsocketError{Code: GENERIC_RECONNECT_IMMEDIATELY, Msg: "Generic reconnect immediately"},
	}
}

// When pusher wants to send an Generic error, it only send the message, the code become nil
// Currently I do not know how to send nil, so I send GENERIC_ERROR
type GenericError struct {
	BaseWebsocketError
}

func newGenericError(msg string) GenericError {
	return GenericError{
		BaseWebsocketError{Code: GENERIC_ERROR, Msg: msg},
	}
}
