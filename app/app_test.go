// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	"strconv"
	"testing"

	channel2 "ipe/channel"
	"ipe/connection"
	"ipe/mocks"
)

var id = 0

func newTestApp() *Application {
	a := NewApplication("Test", strconv.Itoa(id), "123", "123", false, false, true, false, "")
	id++

	return a
}

func TestConnect(t *testing.T) {
	app := newTestApp()

	app.Connect(connection.New("socketID", mocks.MockSocket{}))

	if len(app.connections) != 1 {
		t.Errorf("len(Application.connections) == %d, wants %d", len(app.connections), 1)
	}

}

func TestDisconnect(t *testing.T) {
	app := newTestApp()

	app.Connect(connection.New("socketID", mocks.MockSocket{}))
	app.Disconnect("socketID")

	if len(app.connections) != 0 {
		t.Errorf("len(Application.connections) == %d, wants %d", len(app.connections), 0)
	}

}

func TestFindConnection(t *testing.T) {
	app := newTestApp()

	app.Connect(connection.New("socketID", mocks.MockSocket{}))

	if _, err := app.FindConnection("socketID"); err != nil {
		t.Errorf("Application.FindConnection('socketID') == _, %q, wants %v", err, nil)
	}

	if _, err := app.FindConnection("NotFound"); err == nil {
		t.Errorf("Application.FindConnection('socketID') == _, %q, wants !nil", err)
	}

}

func TestFindChannelByChannelID(t *testing.T) {
	app := newTestApp()

	channel := channel2.New("ID")
	app.AddChannel(channel)

	if _, err := app.FindChannelByChannelID("ID"); err != nil {
		t.Errorf("Application.FindChannelByChannelID('ID') == _, %q, wants %v", err, nil)
	}
}

func TestFindOrCreateChannelByChannelID(t *testing.T) {
	app := newTestApp()

	if len(app.channels) != 0 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 0)
	}

	app.FindOrCreateChannelByChannelID("ID")

	if len(app.channels) != 1 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 1)
	}

}

func TestRemoveChannel(t *testing.T) {
	app := newTestApp()

	if len(app.channels) != 0 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 0)
	}

	channel := channel2.New("ID")
	app.AddChannel(channel)

	if len(app.channels) != 1 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 1)
	}

	app.RemoveChannel(channel)

	if len(app.channels) != 0 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 0)
	}

}

func Test_add_channels(t *testing.T) {

	app := newTestApp()

	// Public

	if len(app.PublicChannels()) != 0 {
		t.Errorf("len(Application.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 0)
	}

	app.AddChannel(channel2.New("ID"))

	if len(app.PublicChannels()) != 1 {
		t.Errorf("len(Application.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 1)
	}

	// Presence

	if len(app.PresenceChannels()) != 0 {
		t.Errorf("len(Application.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 0)
	}

	app.AddChannel(channel2.New("presence-test"))

	if len(app.PresenceChannels()) != 1 {
		t.Errorf("len(Application.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 1)
	}

	// Private

	if len(app.PrivateChannels()) != 0 {
		t.Errorf("len(Application.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 0)
	}

	app.AddChannel(channel2.New("private-test"))

	if len(app.PrivateChannels()) != 1 {
		t.Errorf("len(Application.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 1)
	}

}

func Test_AllChannels(t *testing.T) {
	app := newTestApp()
	app.AddChannel(channel2.New("private-test"))
	app.AddChannel(channel2.New("presence-test"))
	app.AddChannel(channel2.New("test"))

	if len(app.channels) != 3 {
		t.Errorf("len(Application.channels) == %d, wants %d", len(app.channels), 3)
	}
}

func Test_New_Subscriber(t *testing.T) {
	app := newTestApp()

	if len(app.connections) != 0 {
		t.Errorf("len(Application.connections) == %d, wants %d", len(app.connections), 0)
	}

	conn := connection.New("1", mocks.MockSocket{})
	app.Connect(conn)

	if len(app.connections) != 1 {
		t.Errorf("len(Application.connections) == %d, wants %d", len(app.connections), 1)
	}
}

func Test_find_subscriber(t *testing.T) {
	app := newTestApp()
	conn := connection.New("1", mocks.MockSocket{})
	app.Connect(conn)

	conn, err := app.FindConnection("1")

	if err != nil {
		t.Error(err)
	}

	if conn.SocketID != "1" {
		t.Errorf("conn.SocketID == %s, wants %s", conn.SocketID, "1")
	}

	// Find a wrong subscriber

	conn, err = app.FindConnection("DoesNotExists")

	if err == nil {
		t.Errorf("err == %q, wants !nil", err)
	}

	if conn != nil {
		t.Errorf("conn == %q, wants nil", conn)
	}
}

func Test_find_or_create_channels(t *testing.T) {
	app := newTestApp()

	// Public
	if len(app.PublicChannels()) != 0 {
		t.Errorf("len(Application.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 0)
	}

	c := app.FindOrCreateChannelByChannelID("id")

	if len(app.PublicChannels()) != 1 {
		t.Errorf("len(Application.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 1)
	}

	if c.ID != "id" {
		t.Errorf("c.id == %s, wants %s", c.ID, "id")
	}

	// Presence
	if len(app.PresenceChannels()) != 0 {
		t.Errorf("len(Application.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 0)
	}

	c = app.FindOrCreateChannelByChannelID("presence-test")

	if len(app.PresenceChannels()) != 1 {
		t.Errorf("len(Application.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 1)
	}

	if c.ID != "presence-test" {
		t.Errorf("c.id == %s, wants %s", c.ID, "presence-test")
	}

	// Private
	if len(app.PrivateChannels()) != 0 {
		t.Errorf("len(Application.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 0)
	}

	c = app.FindOrCreateChannelByChannelID("private-test")

	if len(app.PrivateChannels()) != 1 {
		t.Errorf("len(Application.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 1)
	}

	if c.ID != "private-test" {
		t.Errorf("c.id == %s, wants %s", c.ID, "private-test")
	}

}
