// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Conf holds the global configuration state
var Conf ConfigFile

// Start Parse the configuration file and starts the ipe server
func Start(configfile string) error {
	file, err := ioutil.ReadFile(configfile)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &Conf); err != nil {
		return err
	}

	Conf.Init()
	router := NewRouter()

	if err := http.ListenAndServe(Conf.Host, router); err != nil {
		return err
	}

	return nil
}
