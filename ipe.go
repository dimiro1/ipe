// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"encoding/json"

	log "github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

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

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return
	}

	// Expand env vars
	data = []byte(os.ExpandEnv(string(data)))

	// Decoding config
	if err := yaml.UnmarshalStrict(data, &conf); err != nil {
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
			a.Enabled,
			a.UserEvents,
			a.WebHooks.Enabled,
			a.WebHooks.URL,
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

	router.Path("/app/create").Methods("POST").Handler(
		api.NewCreateApplication(inMemoryStorage),
	)

	router.HandleFunc("/apps/all",
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != os.Getenv("AUTH_TOKEN") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{ "Error":  "Unauthorized access" }`))
				return
			}

			out, err := json.Marshal(inMemoryStorage.GetAllApps())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ ` + string(out) + ` }`))
			if err != nil {
				log.Error("Error", err)
			}
		},
	)

	router.HandleFunc("/app/test",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`Working Test route`))
		},
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

	if conf.SSL.Enabled {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSL.Host)
			log.Fatal(http.ListenAndServeTLS(conf.SSL.Host, conf.SSL.CertFile, conf.SSL.KeyFile, router))
		}()
	}

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, router))
}
