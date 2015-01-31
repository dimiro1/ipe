// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"strconv"
	"testing"
)

var id = 0

func newApp() *App {

	a := App{Name: "Test", AppID: strconv.Itoa(id), Key: "123", Secret: "123", OnlySSL: false, ApplicationDisabled: false, UserEvents: true}
	a.Init()

	id++
	return &a
}

func TestConnect(t *testing.T) {
	app := newApp()

	app.Connect(newConnection("socketID", nil))

	if len(app.Connections) != 1 {
		t.Errorf("Connections must be 1, but was %d", len(app.Connections))
	}

}

func TestDisconnect(t *testing.T) {
	app := newApp()

	app.Connect(newConnection("socketID", nil))
	app.Disconnect("socketID")

	if len(app.Connections) != 0 {
		t.Errorf("Connections must be 0, but was %d", len(app.Connections))
	}

}

func TestFindConnection(t *testing.T) {
	app := newApp()

	app.Connect(newConnection("socketID", nil))

	if _, err := app.FindConnection("socketID"); err != nil {
		t.Error("Must find Connection")
	}

	if _, err := app.FindConnection("NotFound"); err == nil {
		t.Error("Must not found Connection")
	}

}

func TestFindChannelByChannelID(t *testing.T) {
	app := newApp()

	channel := newChannel("ID")
	app.AddChannel(channel)

	if _, err := app.FindChannelByChannelID("ID"); err != nil {
		t.Error("Channel not found")
	}
}

func TestFindOrCreateChannelByChannelID(t *testing.T) {
	app := newApp()

	if len(app.Channels) != 0 {
		t.Error("Length of channels must be 0 before test")
	}

	app.FindOrCreateChannelByChannelID("ID")

	if len(app.Channels) != 1 {
		t.Error("Length of channels must be 1 after test")
	}

}

func TestRemoveChannel(t *testing.T) {
	app := newApp()

	if len(app.Channels) != 0 {
		t.Error("Length of channels must be 0 before test")
	}

	channel := newChannel("ID")
	app.AddChannel(channel)

	if len(app.Channels) != 1 {
		t.Error("Length of channels after insert must be 1")
	}

	app.RemoveChannel(channel)

	if len(app.Channels) != 0 {
		t.Error("Length of channels must be 0 after remove")
	}

}

func Test_add_channels(t *testing.T) {

	app := newApp()

	// Public

	if len(app.PublicChannels()) != 0 {
		t.Error("Length of public channels must be 0 before test")
	}

	app.AddChannel(newChannel("ID"))

	if len(app.PublicChannels()) != 1 {
		t.Error("Length os public channels after insert must be 1")
	}

	// Presence

	if len(app.PresenceChannels()) != 0 {
		t.Error("Length of presence channels must be 0 before test")
	}

	app.AddChannel(newChannel("presence-test"))

	if len(app.PresenceChannels()) != 1 {
		t.Error("Length os presence channels after insert must be 1")
	}

	// Private

	if len(app.PrivateChannels()) != 0 {
		t.Error("Length of private channels must be 0 before test")
	}

	app.AddChannel(newChannel("private-test"))

	if len(app.PrivateChannels()) != 1 {
		t.Error("Length os private channels after insert must be 1")
	}

}

func Test_AllChannels(t *testing.T) {
	app := newApp()
	app.AddChannel(newChannel("private-test"))
	app.AddChannel(newChannel("presence-test"))
	app.AddChannel(newChannel("test"))

	if len(app.Channels) != 3 {
		t.Error("Must have 3 channels")
	}
}

func Test_New_Subscriber(t *testing.T) {
	app := newApp()

	if len(app.Connections) != 0 {
		t.Error("Length of subscribers before test must be 0")
	}

	conn := newConnection("1", nil)
	app.Connect(conn)

	if len(app.Connections) != 1 {
		t.Error("Length os subscribers after test must be 1")
	}
}

func Test_find_subscriber(t *testing.T) {
	app := newApp()
	conn := newConnection("1", nil)
	app.Connect(conn)

	conn, err := app.FindConnection("1")

	if err != nil {
		t.Error(err)
	}

	if conn.SocketID != "1" {
		t.Error("Wrong subscriber.")
	}

	// Find a wrong subscriber

	conn, err = app.FindConnection("DoesNotExists")

	if err == nil {
		t.Error("Opps, Must be nil")
	}

	if conn != nil {
		t.Error("Opps, Must be nil")
	}
}

func Test_find_or_create_channels(t *testing.T) {
	app := newApp()

	// Public
	if len(app.PublicChannels()) != 0 {
		t.Error("Length of public channels must be 0 before test")
	}

	c := app.FindOrCreateChannelByChannelID("id")

	if len(app.PublicChannels()) != 1 {
		t.Error("Length os public channels after insert must be 1")
	}

	if c.ChannelID != "id" {
		t.Error("Opps wrong channel")
	}

	// Presence
	if len(app.PresenceChannels()) != 0 {
		t.Error("Length of presence channels must be 0 before test")
	}

	c = app.FindOrCreateChannelByChannelID("presence-test")

	if len(app.PresenceChannels()) != 1 {
		t.Error("Length os presence channels after insert must be 1")
	}

	if c.ChannelID != "presence-test" {
		t.Error("Opps wrong channel")
	}

	// Private
	if len(app.PrivateChannels()) != 0 {
		t.Error("Length of private channels must be 0 before test")
	}

	c = app.FindOrCreateChannelByChannelID("private-test")

	if len(app.PrivateChannels()) != 1 {
		t.Error("Length os private channels after insert must be 1")
	}

	if c.ChannelID != "private-test" {
		t.Error("Opps wrong channel")
	}

}
