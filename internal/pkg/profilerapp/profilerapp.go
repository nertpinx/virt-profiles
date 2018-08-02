/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 */

package profilerapp

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	catalogue "github.com/fromanirh/virt-profiles/pkg/catalogue"
	"github.com/gorilla/mux"
)

type ProfilerApp struct {
	cat *catalogue.Catalogue
	mux *mux.Router
}

func NewProfilerApp(profilesDir string) (*ProfilerApp, error) {
	cat, err := catalogue.NewCatalogue(profilesDir)
	if err != nil {
		return nil, err
	}
	app := &ProfilerApp{
		cat: cat,
		mux: mux.NewRouter().StrictSlash(true),
	}
	// POST: receive a preset from KubeVirt, add it to the profiles catalogue
	app.mux.HandleFunc("/presets", app.Presets)
	// GET: list all the profiles known to the system
	app.mux.HandleFunc("/profiles", app.Profiles)
	// POST: apply all the relevant profiles to the domainspec, return updated domainspec
	app.mux.HandleFunc("/domainspec", app.DomainSpec)
	return app, nil
}

func (pa *ProfilerApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pa.mux.ServeHTTP(w, req)
}

type appError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorResponse(w http.ResponseWriter, httpCode, errCode int, errMessage string) {
	w.WriteHeader(httpCode)
	enc := json.NewEncoder(w)
	msg := appError{Code: errCode, Message: errMessage}
	err := enc.Encode(msg)
	if err != nil {
		w.Write([]byte("500 - Something bad happened!"))
	}
}

func (pa *ProfilerApp) Profiles(w http.ResponseWriter, r *http.Request) {
	entries, err := pa.cat.Names()
	if err != nil {
		log.Printf("profiles: gathering: %v", err)
		errorResponse(w, http.StatusInternalServerError, 0, err.Error())
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(entries)
	if err != nil {
		log.Printf("profiles: encoding: %v", err)
		errorResponse(w, http.StatusInternalServerError, 0, err.Error())
		return
	}
}

func (pa *ProfilerApp) Presets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Presets, %q", html.EscapeString(r.URL.Path))
}

func (pa *ProfilerApp) DomainSpec(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DomainSpec, %q", html.EscapeString(r.URL.Path))
}
