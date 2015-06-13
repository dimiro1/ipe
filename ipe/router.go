// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	_ "expvar"
	"net/http"

	"github.com/dimiro1/ipe/vendor/github.com/goji/httpauth"
	"github.com/dimiro1/ipe/vendor/github.com/gorilla/mux"
)

// NewRouter is a function that returns a new configured Router
// It add the necessary middlewares
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	if conf.Expvar {
		if !conf.WasProvidedUserAndPassword() {
			panic("Your are exporting debug variables and looks like you forget to define an User and a Password")
		}

		router.Handle("/debug/vars", httpauth.SimpleBasicAuth(conf.User, conf.Password)(http.DefaultServeMux))
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		if route.RequiresRestAuth {
			handler = restAuthenticationHandler(handler)
			handler = restCheckAppDisabledHandler(handler)
		}

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}
