// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"ipe/api"
	"ipe/app"
	"ipe/config"
	"ipe/storage"
	"ipe/websockets"
)

// Start Parse the configuration file and starts the ipe server
// It Panic if could not start the HTTP or HTTPS server
func Start(filename string) {
	var conf config.File

	rand.Seed(time.Now().Unix())
	file, err := os.Open(filename)

	if err != nil {
		log.Error(err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Reading config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Error(err)
		return
	}

	// Using a in memory database
	inMemoryStorage := storage.NewInMemory()

	// Adding applications
	for _, a := range conf.Apps {
		application := app.NewApplication(
			a.Name,
			a.AppID,
			a.Key,
			a.Secret,
			a.OnlySSL,
			a.ApplicationDisabled,
			a.UserEvents,
			a.WebHooks,
			a.URLWebHook,
		)

		if err := inMemoryStorage.AddApp(application); err != nil {
			log.Error(err)
			return
		}
	}

	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())

	router.Path("/app/{key}").Methods("GET").Handler(
		websockets.NewWebsocket(inMemoryStorage),
	)

	appsRouter := router.PathPrefix("/apps/{app_id}").Subrouter()
	appsRouter.Use(
		api.CheckAppDisabled(inMemoryStorage),
		api.Authentication(inMemoryStorage),
	)

	appsRouter.Path("/events").Methods("POST").Handler(
		api.NewPostEvents(inMemoryStorage),
	)
	appsRouter.Path("/channels").Methods("GET").Handler(
		api.NewGetChannels(inMemoryStorage),
	)
	appsRouter.Path("/channels/{channel_name}").Methods("GET").Handler(
		api.NewGetChannel(inMemoryStorage),
	)
	appsRouter.Path("/channels/{channel_name}/users").Methods("GET").Handler(
		api.NewGetChannelUsers(inMemoryStorage),
	)

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, router))
		}()
	}

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, router))
}
