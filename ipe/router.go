// Copyright 2014, 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"net/http"

	"github.com/gorilla/mux"
)

type router struct {
	ctx    *applicationContext
	mux    *mux.Router
	routes map[string]handlerHTTPC
}

func newRouter(ctx *applicationContext) *router {
	return &router{
		ctx: ctx,
		mux: mux.NewRouter().StrictSlash(true),
	}
}

func (a *router) GET(path string, handler handlerHTTPC) {
	a.mux.Methods("GET").Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPC(a.ctx, w, r)
	})
}

func (a *router) POST(path string, handler handlerHTTPC) {
	a.mux.Methods("POST").Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPC(a.ctx, w, r)
	})
}

func (a router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
