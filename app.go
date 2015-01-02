// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"strings"
)

// Applications
type App struct {
	Name                string
	AppID               string
	Key                 string
	Secret              string
	OnlySSL             bool
	ApplicationDisabled bool

	PublicChannels   []*Channel `json:"-"`
	PresenceChannels []*Channel `json:"-"`
	PrivateChannels  []*Channel `json:"-"`

	Connections []*Connection `json:"-"`
}

func (a *App) AllChannels() []*Channel {
	var channels []*Channel

	for _, c := range a.PrivateChannels {
		channels = append(channels, c)
	}

	for _, c := range a.PublicChannels {
		channels = append(channels, c)
	}

	for _, c := range a.PresenceChannels {
		channels = append(channels, c)
	}

	return channels
}

// Create a new Connection
func (a *App) AddConnection(c *Connection) {
	a.Connections = append(a.Connections, c)
}

func (a *App) FindConnection(socketID string) (*Connection, error) {
	for _, c := range a.Connections {
		if c.SocketID == socketID {
			return c, nil
		}
	}

	return nil, errors.New("Connection not found")
}

func (a *App) AddChannel(c *Channel) {
	if c.isPresence() {
		a.PresenceChannels = append(a.PresenceChannels, c)
	} else if c.isPrivate() {
		a.PrivateChannels = append(a.PrivateChannels, c)
	} else {
		a.PublicChannels = append(a.PublicChannels, c)
	}
}

func (a *App) FindOrCreateChannelByChannelID(n, data string) *Channel {
	channel, err := a.FindChannelByChannelID(n)

	if err != nil {
		channel = NewChannel(n, data)

		a.AddChannel(channel)
	}

	return channel
}

// Find the channel by channel ID
func (a *App) FindChannelByChannelID(n string) (*Channel, error) {
	isPrivate := strings.HasPrefix(n, "private-")
	isPresence := strings.HasPrefix(n, "presence-")

	// Get the channel
	if isPresence {
		for _, channel := range a.PresenceChannels {
			if channel.ChannelID == n {
				return channel, nil
			}
		}
	} else if isPrivate {
		for _, channel := range a.PrivateChannels {
			if channel.ChannelID == n {
				return channel, nil
			}
		}
	} else {
		for _, channel := range a.PublicChannels {
			if channel.ChannelID == n {
				return channel, nil
			}
		}
	}

	return nil, errors.New("Channel does not exists")
}
