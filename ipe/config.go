// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"errors"
	"strings"
)

// The config file
type configFile struct {
	Host     string // The host, eg: :8080 will start on 0.0.0.0:8080
	User     string
	Password string

	// SSL Configurations
	Encrypted     bool
	SSLPrivateKey string
	SSLPublicKey  string

	Apps     []*app
}

// Initialize Apps
func (c *configFile) Init() {
	for _, app := range c.Apps {
		app.Init()
	}
}

func (c *configFile) WasProvidedUserAndPassword() bool {
	return len(strings.TrimSpace(c.User)) > 0 && len(strings.TrimSpace(c.Password)) > 0
}

// Returns an App with by appID
func (c *configFile) GetAppByAppID(appID string) (*app, error) {
	for _, a := range c.Apps {
		if a.AppID == appID {
			return a, nil
		}
	}
	return &app{}, errors.New("App not found")
}

// Returns an App with by key
func (c *configFile) GetAppByKey(key string) (*app, error) {
	for _, a := range c.Apps {
		if a.Key == key {
			return a, nil
		}
	}
	return &app{}, errors.New("App not found")
}
