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

// Conf holds the global configuration state
var conf configFile

// Start Parse the configuration file and starts the ipe server
// It Panic if could not start the HTTP or HTTPS server
func Start(configfile string) {
	rand.Seed(time.Now().Unix())
	file, err := os.Open(configfile)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatal(err)
	}

	conf.Init()
	router := newRouter()

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, router))
		}()
	}

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, router))
}
