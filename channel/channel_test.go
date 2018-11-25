// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package channel

import (
	"ipe/connection"
	"ipe/mocks"
	"ipe/subscription"
	"testing"
)

func TestIsOccupied(t *testing.T) {
	c := New("ID")

	if c.IsOccupied() {
		t.Errorf("c.IsOccupied() == %t, wants %t", c.IsOccupied(), false)
	}

	c.subscriptions["ID"] = subscription.New(connection.New("ID", mocks.MockSocket{}), "")

	if !c.IsOccupied() {
		t.Errorf("c.IsOccupied() == %t, wants %t", c.IsOccupied(), true)
	}
}

func TestIsPrivate(t *testing.T) {
	c := New("private-Channel")

	if !c.IsPrivate() {
		t.Errorf("c.IsPrivate() == %t, wants %t", c.IsPrivate(), true)
	}
}

func TestIsPresence(t *testing.T) {
	c := New("presence-Channel")

	if !c.IsPresence() {
		t.Errorf("c.IsPresence() == %t, wants %t", c.IsPresence(), true)
	}
}

func TestIsPublic(t *testing.T) {
	c := New("Channel")

	if !c.IsPublic() {
		t.Errorf("c.IsPublic() == %t, wants %t", c.IsPublic(), true)
	}
}

func TestIsPrivateOrPresence(t *testing.T) {
	c := New("private-Channel")

	if !c.IsPresenceOrPrivate() {
		t.Errorf("c.IsPresenceOrPrivate() == %t, wants %t", c.IsPresenceOrPrivate(), true)
	}

	c = New("presence-Channel")

	if !c.IsPresenceOrPrivate() {
		t.Errorf("c.IsPresenceOrPrivate() == %t, wants %t", c.IsPresenceOrPrivate(), true)
	}
}

func TestTotalSubscriptions(t *testing.T) {
	c := New("ID")

	if c.TotalSubscriptions() != len(c.subscriptions) {
		t.Errorf("c.TotalSubscriptions() == %d, wants %d", c.TotalSubscriptions(), len(c.subscriptions))
	}
}

func TestTotalUsers(t *testing.T) {
	c := New("ID")

	c.subscriptions["1"] = subscription.New(connection.New("ID", mocks.MockSocket{}), "")
	c.subscriptions["2"] = subscription.New(connection.New("ID", mocks.MockSocket{}), "")

	if c.TotalSubscriptions() != len(c.subscriptions) {
		t.Errorf("c.TotalSubscriptions() == %d, wants %d", c.TotalSubscriptions(), len(c.subscriptions))
	}

	if c.TotalUsers() != 1 {
		t.Errorf("c.TotalUsers() == %d, wants %d", c.TotalUsers(), 1)
	}

}

func TestIsSubscribed(t *testing.T) {
	c := New("ID")
	conn := connection.New("ID", mocks.MockSocket{})

	if c.IsSubscribed(conn) {
		t.Errorf("c.IsSubscribed(%q) == %t, wants %t", conn, c.IsSubscribed(conn), false)
	}

	c.subscriptions["ID"] = subscription.New(conn, "")

	if !c.IsSubscribed(conn) {
		t.Errorf("c.IsSubscribed(%q) == %t, wants %t", conn, c.IsSubscribed(conn), true)
	}
}
