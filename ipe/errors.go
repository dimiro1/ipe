// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

// Base interface
type websocketError interface {
	GetCode() int
	GetMsg() string
}

// Base struct
type baseWebsocketError struct {
	Code int
	Msg  string
}

func (e baseWebsocketError) GetCode() int {
	return e.Code
}

func (e baseWebsocketError) GetMsg() string {
	return e.Msg
}

// Unsupprted protocol version
type unsupportedProtocolVersionError struct {
	baseWebsocketError
}

func newUnsupportedProtocolVersionError() unsupportedProtocolVersionError {
	return unsupportedProtocolVersionError{
		baseWebsocketError{Code: UNSUPPORTED_PROTOCOL_VERSION, Msg: "Unsupported protocol version"},
	}
}

// The application does not exists
// See the configuration file
type applicationDoesNotExistsError struct {
	baseWebsocketError
}

func newApplicationDoesNotExistsError() applicationDoesNotExistsError {
	return applicationDoesNotExistsError{
		baseWebsocketError{Code: APPLICATION_DOES_NOT_EXISTS, Msg: "Could not found an app with the given key"},
	}
}

// The user did not send the protocol version
type noProtocolVersionSuppliedError struct {
	baseWebsocketError
}

func newNoProtocolVersionSuppliedError() noProtocolVersionSuppliedError {
	return noProtocolVersionSuppliedError{
		baseWebsocketError{Code: NO_PROTOCOL_VERSION_SUPPLIED, Msg: "No protocol version supplied"},
	}
}

// When the application is disabled.
// See the configuration file
type applicationDisabledError struct {
	baseWebsocketError
}

func newApplicationDisabledError() noProtocolVersionSuppliedError {
	return noProtocolVersionSuppliedError{
		baseWebsocketError{Code: APPLICATION_DISABLED, Msg: "Application disabled"},
	}
}

// When the application only accepts SSL connections
type applicationOnlyAccepsSSLError struct {
	baseWebsocketError
}

func newApplicationOnlyAccepsSSLError() applicationOnlyAccepsSSLError {
	return applicationOnlyAccepsSSLError{
		baseWebsocketError{Code: APPLICATION_ONLY_ACCEPTS_SSL, Msg: "Application only accepts SSL connections, reconnect using wss://"},
	}
}

// When the user send an invalid version
type invalidVersionStringFormatError struct {
	baseWebsocketError
}

func newInvalidVersionStringFormatError() invalidVersionStringFormatError {
	return invalidVersionStringFormatError{
		baseWebsocketError{Code: INVALID_VERSION_STRING_FORMAT, Msg: "Invalid version string format"},
	}
}

// Used when the error was internal
// * Decoding json
// * Writing to output
type genericReconnectImmediatelyError struct {
	baseWebsocketError
}

func newGenericReconnectImmediatelyError() genericReconnectImmediatelyError {
	return genericReconnectImmediatelyError{
		baseWebsocketError{Code: GENERIC_RECONNECT_IMMEDIATELY, Msg: "Generic reconnect immediately"},
	}
}

// When pusher wants to send an Generic error, it only send the message, the code become nil
// Currently I do not know how to send nil, so I send GENERIC_ERROR
type genericError struct {
	baseWebsocketError
}

func newGenericError(msg string) genericError {
	return genericError{
		baseWebsocketError{Code: GENERIC_ERROR, Msg: msg},
	}
}
