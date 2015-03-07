// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
)

// Conf holds the global configuration state
var conf configFile

// Start Parse the configuration file and starts the ipe server
func Start(configfile string) error {
	file, err := ioutil.ReadFile(configfile)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &conf); err != nil {
		return err
	}

	conf.Init()
	if strings.HasPrefix(conf.Host, ":$") {
    conf.Host = ":" + os.Getenv(conf.Host[2:len(conf.Host)])
  }
	router := newRouter()

	if err := http.ListenAndServe(conf.Host, router); err != nil {
		return err
	}

	return nil
}
