// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "testing"

func TestIsOccupied(t *testing.T) {
	c := newChannel("ID")

	if c.IsOccupied() {
		t.Errorf("c.IsOccupied() == %t, wants %t", c.IsOccupied(), false)
	}

	c.Subscriptions["ID"] = newSubscription(newConnection("ID", nil), "")

	if !c.IsOccupied() {
		t.Errorf("c.IsOccupied() == %t, wants %t", c.IsOccupied(), true)
	}
}

func TestIsPrivate(t *testing.T) {
	c := newChannel("private-channel")

	if !c.IsPrivate() {
		t.Errorf("c.IsPrivate() == %t, wants %t", c.IsPrivate(), true)
	}
}

func TestIsPresence(t *testing.T) {
	c := newChannel("presence-channel")

	if !c.IsPresence() {
		t.Errorf("c.IsPresence() == %t, wants %t", c.IsPresence(), true)
	}
}

func TestIsPublic(t *testing.T) {
	c := newChannel("channel")

	if !c.IsPublic() {
		t.Errorf("c.IsPublic() == %t, wants %t", c.IsPublic(), true)
	}
}

func TestIsPrivateOrPresence(t *testing.T) {
	c := newChannel("private-channel")

	if !c.IsPresenceOrPrivate() {
		t.Errorf("c.IsPresenceOrPrivate() == %t, wants %t", c.IsPresenceOrPrivate(), true)
	}

	c = newChannel("presence-channel")

	if !c.IsPresenceOrPrivate() {
		t.Errorf("c.IsPresenceOrPrivate() == %t, wants %t", c.IsPresenceOrPrivate(), true)
	}
}

func TestTotalSubscriptions(t *testing.T) {
	c := newChannel("ID")

	if c.TotalSubscriptions() != len(c.Subscriptions) {
		t.Errorf("c.TotalSubscriptions() == %d, wants %d", c.TotalSubscriptions(), len(c.Subscriptions))
	}
}

func TestTotalUsers(t *testing.T) {
	c := newChannel("ID")

	c.Subscriptions["1"] = newSubscription(newConnection("ID", nil), "")
	c.Subscriptions["2"] = newSubscription(newConnection("ID", nil), "")

	if c.TotalSubscriptions() != len(c.Subscriptions) {
		t.Errorf("c.TotalSubscriptions() == %d, wants %d", c.TotalSubscriptions(), len(c.Subscriptions))
	}

	if c.TotalUsers() != 1 {
		t.Errorf("c.TotalUsers() == %d, wants %d", c.TotalUsers(), 1)
	}

}

func TestIsSubscribed(t *testing.T) {
	c := newChannel("ID")
	conn := newConnection("ID", nil)

	if c.IsSubscribed(conn) {
		t.Errorf("c.IsSubscribed(%q) == %t, wants %t", conn, c.IsSubscribed(conn), false)
	}

	c.Subscriptions["ID"] = newSubscription(conn, "")

	if !c.IsSubscribed(conn) {
		t.Errorf("c.IsSubscribed(%q) == %t, wants %t", conn, c.IsSubscribed(conn), true)
	}
}
