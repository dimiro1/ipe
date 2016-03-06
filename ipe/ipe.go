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

	// Creating the global application context
	ctx := &applicationContext{DB: db}

	// The router
	router := newRouter(ctx)

	router.POST("/apps/{app_id}/events", commonHandlers(ctx, postEvents))

	router.GET("/apps/{app_id}/channels", commonHandlers(ctx, getChannels))

	router.GET("/apps/{app_id}/channels/{channel_name}", commonHandlers(ctx, getChannel))

	router.GET("/apps/{app_id}/channels/{channel_name}/users", commonHandlers(ctx, getChannelUsers))

	router.GET("/app/{key}", handlerHTTPCFunc(wsHandler))

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, router))
		}()
	}

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, router))
}
