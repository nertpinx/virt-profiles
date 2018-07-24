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

func (pa *ProfilesApp) Profiles(w http.ResponseWriter, r *http.Request) {
	entries, err := pa.cat.Names()
	if err != nil {
		log.Printf("profiles: gathering: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(entries)
	if err != nil {
		log.Printf("profiles: encoding: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
}

func (pa *ProfilesApp) Resources(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Resources, %q", html.EscapeString(r.URL.Path))
}

func (pa *ProfilesApp) DomainSpec(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DomainSpec, %q", html.EscapeString(r.URL.Path))
}
