// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package subscription

import "ipe/connection"

// Subscription A Channel Subscription
type Subscription struct {
	Connection *connection.Connection
	ID         string
	Data       string
}

// New Create a new Subscription
func New(conn *connection.Connection, data string) *Subscription {
	return &Subscription{Connection: conn, Data: data}
}
