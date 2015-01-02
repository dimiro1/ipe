// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var Conf ConfigFile

func main() {
	var filename = flag.String("config", "config.json", "Config file location")

	flag.Parse()

	file, err := ioutil.ReadFile(*filename)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &Conf); err != nil {
		log.Fatal(err)
	}

	router := NewRouter()

	log.Printf("Starting Ipê using config file: '%s'", *filename)

	if err := http.ListenAndServe(Conf.Host, router); err != nil {
		log.Fatalln(err)
	}
}
