// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// Start Parse the configuration file and starts the ipe server
// It Panic if could not start the HTTP or HTTPS server
func Start(filename string) {
	conf := loadConfig(filename)

	// Using a in memory database
	db := newMemdb()

	// Adding applications
	for _, a := range conf.Apps {
		db.AddApp(newAppFromConfig(a))
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/app/:key", (&websocketHandler{db}).ServeHTTP)
	r.Group(func(r chi.Router) {
		r.Use(checkAppDisabled(db))
		r.Use(authenticationHandler(db))

		r.Post("/apps/:app_id/events", (&postEventsHandler{db}).ServeHTTP)
		r.Get("/apps/:app_id/channels", (&getChannelsHandler{db}).ServeHTTP)
		r.Get("/apps/:app_id/channels/:channel_name", (&getChannelHandler{db}).ServeHTTP)
		r.Get("/apps/:app_id/channels/:channel_name/users", (&getChannelUsersHandler{db}).ServeHTTP)
	})

	if conf.Profiling {
		r.Mount("/debug", middleware.Profiler())
	}

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, r))
		}()
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)
	go func() {
		for {
			<-s
			conf = loadConfig(filename)
			for _, a := range conf.Apps {
				if exists, _ := db.GetAppByAppID(a.AppID); exists != nil {
					continue
				}
				db.AddApp(newAppFromConfig(a))
			}
			log.Info("Reloaded config")
		}
	}()

	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, r))
}
