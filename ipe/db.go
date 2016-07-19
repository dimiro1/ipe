// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"errors"
	"sync"
)

// db represents a app database
// For now it there is only one memory database implementation
// but in the future I can write a sql implementation
type db interface {
	GetAppByAppID(appID string) (*app, error)
	GetAppByKey(key string) (*app, error)
	AddApp(*app) error
}

// memdb is a in memory implementation of db interface
type memdb struct {
	IdMutex     sync.Mutex
	KeyMutex    sync.Mutex
	AppsByAppID map[string]*app
	AppsByKey   map[string]*app
}

func newMemdb() *memdb {
	return &memdb{}
}

func (db *memdb) AddApp(a *app) error {
	db.IdMutex.Lock()
	db.AppsByAppID[a.AppID] = a
	db.IdMutex.Unlock()

	db.KeyMutex.Lock()
	db.AppsByKey[a.Key] = a
	db.KeyMutex.Unlock()

	return nil
}

// GetAppByAppID returns an App with by appID
func (db *memdb) GetAppByAppID(appID string) (*app, error) {
	db.IdMutex.Lock()
	a, exists := db.AppsByAppID[appID]
	db.IdMutex.Unlock()
	if exists {
		return a, nil
	}
	return nil, errors.New("App not found")
}

// GetAppByKey returns an App with by key
func (db *memdb) GetAppByKey(key string) (*app, error) {
	db.KeyMutex.Lock()
	a, exists := db.AppsByKey[key]
	db.KeyMutex.Unlock()
	if exists {
		return a, nil
	}
	return nil, errors.New("App not found")
}
