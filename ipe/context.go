// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "net/http"

type applicationContext struct {
	DB db
}

// url params
type params map[string]string

func (p params) Get(key string) string {
	return p[key]
}

// A contextHandler responds to an HTTP request with custom application context.
type contextHandler interface {
	ServeWithContext(ctx *applicationContext, p params, w http.ResponseWriter, r *http.Request)
}

type contextHandlerFunc func(ctx *applicationContext, p params, w http.ResponseWriter, r *http.Request)

func (c contextHandlerFunc) ServeWithContext(ctx *applicationContext, p params, w http.ResponseWriter, r *http.Request) {
	c(ctx, p, w, r)
}
