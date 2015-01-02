// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
)

type ConfigFile struct {
	Host          string
	SessionName   string
	SessionSecret string
	Apps          []*App
}

func (c ConfigFile) GetAppByAppID(appID string) (*App, error) {
	for _, app := range c.Apps {
		if app.AppID == appID {
			return app, nil
		}
	}
	return &App{}, errors.New("App not found")
}

func (c ConfigFile) GetAppByKey(key string) (*App, error) {
	for _, app := range c.Apps {
		if app.Key == key {
			return app, nil
		}
	}
	return &App{}, errors.New("App not found")
}
