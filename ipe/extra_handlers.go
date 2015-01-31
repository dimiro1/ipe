// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Check if the application is disabled
func RestCheckAppDisabledHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appID := vars["app_id"]

		currentApp, err := Conf.GetAppByAppID(appID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Could not found an app with app_id: %s", appID), http.StatusForbidden)
			return
		}

		if currentApp.ApplicationDisabled {
			http.Error(w, "Application disabled", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
