package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/fromanirh/virt-profiles/pkg/virtprofiles"
	"github.com/gorilla/mux"
)

type ProfilesApp struct {
	cat *virtprofiles.Catalogue
	mux *mux.Router
}

func NewProfilesApp(profilesDir string) (*ProfilesApp, error) {
	cat, err := virtprofiles.NewCatalogue(profilesDir)
	if err != nil {
		return nil, err
	}
	app := &ProfilesApp{
		cat: cat,
		mux: mux.NewRouter().StrictSlash(true),
	}
	app.mux.HandleFunc("/profiles", app.Profiles)
	app.mux.HandleFunc("/resources", app.Resources)
	app.mux.HandleFunc("/domainspec", app.DomainSpec)
	return app, nil
}

func (pa *ProfilesApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

func (pa *ProfilesApp) Profiles(w http.ResponseWriter, r *http.Request) {
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

func (pa *ProfilesApp) Resources(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Resources, %q", html.EscapeString(r.URL.Path))
}

func (pa *ProfilesApp) DomainSpec(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DomainSpec, %q", html.EscapeString(r.URL.Path))
}
