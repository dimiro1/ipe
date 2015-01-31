// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

// A Channel Subscription
type subscription struct {
	Connection *connection
	Id         string
	Data       string
}

// Create a new Subscription
func newSubscription(conn *connection, data string) *subscription {
	return &subscription{Connection: conn, Data: data}
}
