// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"net/http"

	"github.com/gorilla/mux"
)

// newRouter is a function that returns a new configured Router
// It add the necessary middlewares
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

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
