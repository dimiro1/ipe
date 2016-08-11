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

	"goji.io/pat"

	goji "goji.io"

	log "github.com/golang/glog"
)

// Start Parse the configuration file and starts the ipe server
// It Panic if could not start the HTTP or HTTPS server
func Start(filename string) {
	var conf configFile

	rand.Seed(time.Now().Unix())
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Reading config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatal(err)
	}

	// Using a in memory database
	db := newMemdb()

	// Adding applications
	for _, a := range conf.Apps {
		db.AddApp(newAppFromConfig(a))
	}

	router := goji.NewMux()
	router.HandleFuncC(pat.Post("/apps/:app_id/events"), newPostEventsHandler(db))
	router.HandleFuncC(pat.Get("/apps/:app_id/channels"), newGetChannelsHandler(db))
	router.HandleFuncC(pat.Get("/apps/:app_id/channels/:channel_name"), newGetChannelHandler(db))
	router.HandleFuncC(pat.Get("/apps/:app_id/channels/:channel_name/users"), newGetChannelUsersHandler(db))
	router.HandleC(pat.Get("/app/:key"), newWebsocketHandler(db))

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, router))
		}()
	}

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, router))
}
