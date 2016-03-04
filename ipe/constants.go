// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

// Error Codes
const (
	// 4000 - 4099
	// Indicates an error resulting in the connection being closed by Pusher,
	// and that attempting to reconnect using the same parameters will not succeed.
	applicationOnlyAcceptsSSL        = 4000
	applicationDoesNotExists         = 4001
	applicationDisabled              = 4003
	applicationIsOverConnectionQuota = 4004 // Not Implemented
	pathNotFound                     = 4005 // Not Implemented
	invalidVersionStringFormat       = 4006
	unsupportedProtocolVersion       = 4007
	noProtocolVersionSupplied        = 4008

	// 4100 - 4199
	// Indicates an error resulting in the connection being closed by Pusher,
	// and the client may reconnect after 1s or more
	overCapacity = 4100 // Not Implemented

	// 4200 - 4299
	// Indicate an error resulting in the connection being closed by Pusher,
	// and the client my reconnect immediately
	genericReconnectImmediately = 4200
	pongReplyNotReceived        = 4201 // Ping was sent to the client, but no reply was received; Not Implemented
	closedAfterInactivity       = 4202 // Client has been inactive for a long time (24 hours) and client does not suppot ping.; Not Implemented

	// 4300 - 4399
	// Any other type of error
	clientRejectedDueToRateLimit = 4301 // Not Implemented

	// Pusher send null, This app use this error code to send the null value
	// see ErrorEvent
	otherError = 0
)

// Only this version is supported
const supportedProtocolVersion = 7

// // Maximun event size permitted 10 kB
// See: http://blogs.gnome.org/cneumair/2008/09/30/1-kb-1024-bytes-no-1-kb-1000-bytes/
const maxDataEventSize = 10 * 1000
