// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "testing"

func TestIsOccupied(t *testing.T) {
	c := NewChannel("ID")

	if c.IsOccupied() {
		t.Error("Channels must be empty")
	}

	c.Subscriptions["ID"] = newSubscription(newConnection("ID", nil), "")

	if !c.IsOccupied() {
		t.Error("Channels must be empty")
	}
}

func TestIsPrivate(t *testing.T) {
	c := NewChannel("private-channel")

	if !c.IsPrivate() {
		t.Error("The Channel must be private")
	}
}

func TestIsPresence(t *testing.T) {
	c := NewChannel("presence-channel")

	if !c.IsPresence() {
		t.Error("The Channel must be presence")
	}
}

func TestIsPublic(t *testing.T) {
	c := NewChannel("channel")

	if !c.IsPublic() {
		t.Error("The Channel must be public")
	}
}

func TestIsPrivateOrPresence(t *testing.T) {
	c := NewChannel("private-channel")

	if !c.IsPresenceOrPrivate() {
		t.Error("The Channel must be private or presence")
	}

	c = NewChannel("presence-channel")

	if !c.IsPresenceOrPrivate() {
		t.Error("The Channel must be private or presence")
	}
}

func TestTotalSubscriptions(t *testing.T) {
	c := NewChannel("ID")

	if c.TotalSubscriptions() != len(c.Subscriptions) {
		t.Error("TotalSubscriptions must be equal to len of total subscriptions")
	}
}

func TestTotalUsers(t *testing.T) {
	c := NewChannel("ID")

	c.Subscriptions["1"] = newSubscription(newConnection("ID", nil), "")
	c.Subscriptions["2"] = newSubscription(newConnection("ID", nil), "")

	if c.TotalSubscriptions() != len(c.Subscriptions) {
		t.Error("TotalSubscriptions must be equal to len of total subscriptions")
	}

	if c.TotalUsers() != 1 {
		t.Error("TotalUsers must be equal to 1")
	}

}

func TestIsSubscribed(t *testing.T) {
	c := NewChannel("ID")
	conn := newConnection("ID", nil)

	if c.IsSubscribed(conn) {
		t.Error("Must not be subscribed")
	}

	c.Subscriptions["ID"] = newSubscription(conn, "")

	if !c.IsSubscribed(conn) {
		t.Error("Must be subscribed")
	}
}
