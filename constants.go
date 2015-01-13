// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import ()

// Error Codes
const (
	// 4000 - 4099
	// Indicates an error resulting in the connection being closed by Pusher,
	// and that attempting to reconnect using the same parameters will not succeed.
	APPLICATION_ONLY_ACCEPTS_SSL         = 4000
	APPLICATION_DOES_NOT_EXISTS          = 4001
	APPLICATION_DISABLED                 = 4003
	APPLICATION_IS_OVER_CONNECTION_QUOTA = 4004 // Not Implemented
	PATH_NOT_FOUND                       = 4005 // Not Implemented
	INVALID_VERSION_STRING_FORMAT        = 4006
	UNSUPPORTED_PROTOCOL_VERSION         = 4007
	NO_PROTOCOL_VERSION_SUPPLIED         = 4008

	// 4100 - 4199
	// Indicates an error resulting in the connection being closed by Pusher,
	// and the client may reconnect after 1s or more
	OVER_CAPACITY = 4100 // Not Implemented

	// 4200 - 4299
	// Indicate an error resulting in the connection being closed by Pusher,
	// and the client my reconnect immediately
	GENERIC_RECONNECT_IMMEDIATELY = 4200
	PONG_REPLY_NOT_RECEIVED       = 4201 // Ping was sent to the client, but no reply was received
	CLOSED_AFTER_INACTIVITY       = 4202 // Client has been inactive for a long time (24 hours) and client does not suppot ping.

	// 4300 - 4399
	// Any other type of error
	CLIENT_REJECTED_DUE_TO_RATE_LIMIT = 4301 // Not Implemented

	// Pusher send null, This app use this error code to send the null value
	// see ErrorEvent
	GENERIC_ERROR = 0
)

// Only this version is supported
const SUPPORTED_PROTOCOL_VERSION = 7

// // Maximun event size permitted 20 kB
// See: http://blogs.gnome.org/cneumair/2008/09/30/1-kb-1024-bytes-no-1-kb-1000-bytes/
const MAX_DATA_EVENT_SIZE = 10 * 1000
