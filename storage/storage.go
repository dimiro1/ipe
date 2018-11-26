// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"errors"
	"ipe/app"
	"sync"
)

// storage represents a app database
// For now it there is only one memory database implementation
// but in the future I can write a sql implementation
type Storage interface {
	GetAppByAppID(appID string) (*app.Application, error)
	GetAppByKey(key string) (*app.Application, error)
	AddApp(application *app.Application) error
}

// InMemory in memory implementation of Storage
type InMemory struct {
	sync.RWMutex
	Apps []*app.Application
}

// NewInMemory returns an InMemory storage
func NewInMemory() Storage {
	return &InMemory{}
}

// AddApp adds app into memory
func (db *InMemory) AddApp(application *app.Application) error {
	db.Lock()
	defer db.Unlock()

	db.Apps = append(db.Apps, application)
	return nil
}

// GetAppByAppID returns an App with by appID
func (db *InMemory) GetAppByAppID(appID string) (*app.Application, error) {
	db.RLock()
	defer db.RUnlock()

	for _, a := range db.Apps {
		if a.AppID == appID {
			return a, nil
		}
	}

	return nil, errors.New("app not found")
}

// GetAppByKey returns an App with by key
func (db *InMemory) GetAppByKey(key string) (*app.Application, error) {
	db.RLock()
	defer db.RUnlock()

	for _, a := range db.Apps {
		if a.Key == key {
			return a, nil
		}
	}
	return nil, errors.New("app not found")
}
