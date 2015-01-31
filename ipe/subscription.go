// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

// A Channel Subscription
type Subscription struct {
	Connection *Connection
	Id         string
	Data       string
}

// Create a new Subscription
func NewSubscription(conn *Connection, data string) *Subscription {
	return &Subscription{Connection: conn, Data: data}
}
