// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/golang/glog"
)

var Conf ConfigFile

func main() {
	printBanner()

	var filename = flag.String("config", "config.json", "Config file location")

	flag.Parse()

	file, err := ioutil.ReadFile(*filename)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &Conf); err != nil {
		log.Fatal(err)
	}

	Conf.Init()

	router := NewRouter()

	log.Infof("Starting Ipê using config file: '%s'", *filename)

	if err := http.ListenAndServe(Conf.Host, router); err != nil {
		log.Fatal(err)
	}
}

func printBanner() {
	fmt.Println("\033[32mWelcome to Ipê - Yet another Pusher server clone\033[0m")
	fmt.Println("\033[33mBy: Claudemiro Alves Feitosa Neto <dimiro1@gmail.com>\033[0m")
}
