// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"math/rand"
	"time"

	"github.com/gorilla/mux"
	log "github.com/golang/glog"
)

// Conf holds the global configuration state
var conf configFile

func Run(conf configFile, router *mux.Router) chan error {

    errs := make(chan error)

    // Starting HTTP server
    go func() {
        log.Infof("Staring HTTP service on %s ...", conf.Host)

        if err := http.ListenAndServe(conf.Host, router); err != nil {
            errs <- err
        }

    }()

	if conf.Encrypted {
	    // Starting HTTPS server
	    go func() {
	        log.Infof("Staring HTTPS service on %s ...", conf.SSLHost)
	        if err := http.ListenAndServeTLS(conf.SSLHost, conf.SSLPublicKey, conf.SSLPrivateKey, router); err != nil {
	            errs <- err
	        }
	    }()
	}

    return errs
}

// Start Parse the configuration file and starts the ipe server
func Start(configfile string) {
	rand.Seed(time.Now().Unix())
	file, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &conf); err != nil {
		log.Fatal(err)
	}

	conf.Init()
	router := newRouter()

	errs := Run(conf, router)

	select {
    case err := <-errs:
        log.Errorf("Could not start serving service due to (error: %s)", err)
    }
}
