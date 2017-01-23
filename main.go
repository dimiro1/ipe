// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	"github.com/rtatibouet/ipe/ipe"
)

// These variables are generated by the linker
// please see the makefile for mor information.
var (
	version    = "version"
	buildstamp = "buildstamp"
	githash    = "githash"
)

// Main function, initialize the system
func main() {
	var filename = flag.String("config", "config.json", "Config file location")
	flag.Parse()

	printBanner()

	ipe.Start(*filename)
}

// Print a beautiful banner
func printBanner() {
	fmt.Print("\033[31m")
	fmt.Print(`  
d8b             db     
Y8P           d8  8b        
                      
888 88888b.   .d88b.  
888 888 "88b d8P  Y8b 
888 888  888 88888888 
888 888 d88P Y8b.     
888 88888P"   "Y8888  
    888               
    888               
    888               
`)
	fmt.Println("\033[0m")
	fmt.Println("\033[32mWelcome to Ipê - Yet another Pusher server clone (https://github.com/dimiro1/ipe)\033[0m")
	fmt.Printf("\033[32mVersion %s+%s.git.%s\033[0m\n", version, buildstamp, githash)
	fmt.Println("\033[33mBy: Claudemiro Alves Feitosa Neto <dimiro1@gmail.com>\033[0m")
}
