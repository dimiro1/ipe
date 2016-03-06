// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "net/http"

type applicationContext struct {
	DB db
}

// A handlerHTTPC responds to an HTTP request with custom application context.
type handlerHTTPC interface {
	ServeHTTPC(ctx *applicationContext, w http.ResponseWriter, r *http.Request)
}

type handlerHTTPCFunc func(ctx *applicationContext, w http.ResponseWriter, r *http.Request)

func (c handlerHTTPCFunc) ServeHTTPC(ctx *applicationContext, w http.ResponseWriter, r *http.Request) {
	c(ctx, w, r)
}
