// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func newApp() App {
	return App{Name: "Test", AppID: "123", Key: "123", Secret: "123", OnlySSL: false, ApplicationDisabled: false}
}

func Test_add_channels(t *testing.T) {

	app := newApp()

	// Public

	if len(app.PublicChannels) != 0 {
		t.Error("Length of public channels must be 0 before test")
	}

	app.AddChannel(NewChannel("ID", ""))

	if len(app.PublicChannels) != 1 {
		t.Error("Length os public channels after insert must be 1")
	}

	// Presence

	if len(app.PresenceChannels) != 0 {
		t.Error("Length of presence channels must be 0 before test")
	}

	app.AddChannel(NewChannel("presence-test", ""))

	if len(app.PresenceChannels) != 1 {
		t.Error("Length os presence channels after insert must be 1")
	}

	// Private

	if len(app.PrivateChannels) != 0 {
		t.Error("Length of private channels must be 0 before test")
	}

	app.AddChannel(NewChannel("private-test", ""))

	if len(app.PrivateChannels) != 1 {
		t.Error("Length os private channels after insert must be 1")
	}

}

func Test_AllChannels(t *testing.T) {
	app := newApp()
	app.AddChannel(NewChannel("private-test", ""))
	app.AddChannel(NewChannel("presence-test", ""))
	app.AddChannel(NewChannel("test", ""))

	if len(app.AllChannels()) != 3 {
		t.Error("Must have 3 channels")
	}
}

func Test_New_Connection(t *testing.T) {
	app := newApp()

	if len(app.Connections) != 0 {
		t.Error("Length of connections before test must be 0")
	}

	conn := NewConnection("1", "", nil)
	app.AddConnection(conn)

	if len(app.Connections) != 1 {
		t.Error("Length os connections after test must be 1")
	}
}

func Test_find_connection(t *testing.T) {
	app := newApp()
	conn := NewConnection("1", "", nil)
	app.AddConnection(conn)

	conn, err := app.FindConnection("1")

	if err != nil {
		t.Error(err)
	}

	if conn.SocketID != "1" {
		t.Error("Wrong connection.")
	}

	// Find a wrong connection

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
	if len(app.PublicChannels) != 0 {
		t.Error("Length of public channels must be 0 before test")
	}

	c := app.FindOrCreateChannelByChannelID("id", "")

	if len(app.PublicChannels) != 1 {
		t.Error("Length os public channels after insert must be 1")
	}

	if c.ChannelID != "id" {
		t.Error("Opps wrong channel")
	}

	// Presence
	if len(app.PresenceChannels) != 0 {
		t.Error("Length of presence channels must be 0 before test")
	}

	c = app.FindOrCreateChannelByChannelID("presence-test", "")

	if len(app.PresenceChannels) != 1 {
		t.Error("Length os presence channels after insert must be 1")
	}

	if c.ChannelID != "presence-test" {
		t.Error("Opps wrong channel")
	}

	// Private
	if len(app.PrivateChannels) != 0 {
		t.Error("Length of private channels must be 0 before test")
	}

	c = app.FindOrCreateChannelByChannelID("private-test", "")

	if len(app.PrivateChannels) != 1 {
		t.Error("Length os private channels after insert must be 1")
	}

	if c.ChannelID != "private-test" {
		t.Error("Opps wrong channel")
	}

}
