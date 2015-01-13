// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import "errors"

// The config file
type ConfigFile struct {
	Host string // The host, eg: :8080 will start on 0.0.0.0:8080
	Apps []*App
}

// Error for App not found
var AppNotFoundError = errors.New("App not found")

func (c *ConfigFile) Initialize() {
	for _, app := range c.Apps {
		app.Subscribers = make(map[string]*Subscriber)
	}
}

// Returns an App with by appID
func (c *ConfigFile) GetAppByAppID(appID string) (*App, error) {
	for _, app := range c.Apps {
		if app.AppID == appID {
			return app, nil
		}
	}
	return &App{}, AppNotFoundError
}

// Returns an App with by key
func (c *ConfigFile) GetAppByKey(key string) (*App, error) {
	for _, app := range c.Apps {
		if app.Key == key {
			return app, nil
		}
	}
	return &App{}, AppNotFoundError
}
