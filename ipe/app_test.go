// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"strconv"
	"testing"
)

var id = 0

func newTestApp() *app {

	a := newApp("Test", strconv.Itoa(id), "123", "123", false, false, true, false, "")
	id++

	return a
}

func TestConnect(t *testing.T) {
	app := newTestApp()

	app.Connect(newConnection("socketID", nil))

	if len(app.Connections) != 1 {
		t.Errorf("len(app.Connections) == %d, wants %d", len(app.Connections), 1)
	}

}

func TestDisconnect(t *testing.T) {
	app := newTestApp()

	app.Connect(newConnection("socketID", nil))
	app.Disconnect("socketID")

	if len(app.Connections) != 0 {
		t.Errorf("len(app.Connections) == %d, wants %d", len(app.Connections), 0)
	}

}

func TestFindConnection(t *testing.T) {
	app := newTestApp()

	app.Connect(newConnection("socketID", nil))

	if _, err := app.FindConnection("socketID"); err != nil {
		t.Errorf("app.FindConnection('socketID') == _, %q, wants %v", err, nil)
	}

	if _, err := app.FindConnection("NotFound"); err == nil {
		t.Errorf("app.FindConnection('socketID') == _, %q, wants !nil", err)
	}

}

func TestFindChannelByChannelID(t *testing.T) {
	app := newTestApp()

	channel := newChannel("ID")
	app.AddChannel(channel)

	if _, err := app.FindChannelByChannelID("ID"); err != nil {
		t.Errorf("app.FindChannelByChannelID('ID') == _, %q, wants %v", err, nil)
	}
}

func TestFindOrCreateChannelByChannelID(t *testing.T) {
	app := newTestApp()

	if len(app.Channels) != 0 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 0)
	}

	app.FindOrCreateChannelByChannelID("ID")

	if len(app.Channels) != 1 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 1)
	}

}

func TestRemoveChannel(t *testing.T) {
	app := newTestApp()

	if len(app.Channels) != 0 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 0)
	}

	channel := newChannel("ID")
	app.AddChannel(channel)

	if len(app.Channels) != 1 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 1)
	}

	app.RemoveChannel(channel)

	if len(app.Channels) != 0 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 0)
	}

}

func Test_add_channels(t *testing.T) {

	app := newTestApp()

	// Public

	if len(app.PublicChannels()) != 0 {
		t.Errorf("len(app.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 0)
	}

	app.AddChannel(newChannel("ID"))

	if len(app.PublicChannels()) != 1 {
		t.Errorf("len(app.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 1)
	}

	// Presence

	if len(app.PresenceChannels()) != 0 {
		t.Errorf("len(app.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 0)
	}

	app.AddChannel(newChannel("presence-test"))

	if len(app.PresenceChannels()) != 1 {
		t.Errorf("len(app.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 1)
	}

	// Private

	if len(app.PrivateChannels()) != 0 {
		t.Errorf("len(app.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 0)
	}

	app.AddChannel(newChannel("private-test"))

	if len(app.PrivateChannels()) != 1 {
		t.Errorf("len(app.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 1)
	}

}

func Test_AllChannels(t *testing.T) {
	app := newTestApp()
	app.AddChannel(newChannel("private-test"))
	app.AddChannel(newChannel("presence-test"))
	app.AddChannel(newChannel("test"))

	if len(app.Channels) != 3 {
		t.Errorf("len(app.Channels) == %d, wants %d", len(app.Channels), 3)
	}
}

func Test_New_Subscriber(t *testing.T) {
	app := newTestApp()

	if len(app.Connections) != 0 {
		t.Errorf("len(app.Connections) == %d, wants %d", len(app.Connections), 0)
	}

	conn := newConnection("1", nil)
	app.Connect(conn)

	if len(app.Connections) != 1 {
		t.Errorf("len(app.Connections) == %d, wants %d", len(app.Connections), 1)
	}
}

func Test_find_subscriber(t *testing.T) {
	app := newTestApp()
	conn := newConnection("1", nil)
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
		t.Errorf("len(app.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 0)
	}

	c := app.FindOrCreateChannelByChannelID("id")

	if len(app.PublicChannels()) != 1 {
		t.Errorf("len(app.PublicChannels()) == %d, wants %d", len(app.PublicChannels()), 1)
	}

	if c.ChannelID != "id" {
		t.Errorf("c.ChannelID == %s, wants %s", c.ChannelID, "id")
	}

	// Presence
	if len(app.PresenceChannels()) != 0 {
		t.Errorf("len(app.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 0)
	}

	c = app.FindOrCreateChannelByChannelID("presence-test")

	if len(app.PresenceChannels()) != 1 {
		t.Errorf("len(app.PresenceChannels()) == %d, wants %d", len(app.PresenceChannels()), 1)
	}

	if c.ChannelID != "presence-test" {
		t.Errorf("c.ChannelID == %s, wants %s", c.ChannelID, "presence-test")
	}

	// Private
	if len(app.PrivateChannels()) != 0 {
		t.Errorf("len(app.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 0)
	}

	c = app.FindOrCreateChannelByChannelID("private-test")

	if len(app.PrivateChannels()) != 1 {
		t.Errorf("len(app.PrivateChannels()) == %d, wants %d", len(app.PrivateChannels()), 1)
	}

	if c.ChannelID != "private-test" {
		t.Errorf("c.ChannelID == %s, wants %s", c.ChannelID, "private-test")
	}

}
