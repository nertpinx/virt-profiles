package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/fromanirh/virt-profiles/pkg/virtprofiles"
	"github.com/gorilla/mux"
)

type ProfilesApp struct {
	cat *virtprofiles.Catalogue
	mux *mux.Router
}

func NewProfilesApp(profilesDir string) (*ProfilesApp, error) {
	app := &ProfilesApp{
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
	fmt.Fprintf(w, "Profiles, %q", html.EscapeString(r.URL.Path))
}

func (pa *ProfilesApp) Resources(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Resources, %q", html.EscapeString(r.URL.Path))
}

func (pa *ProfilesApp) DomainSpec(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DomainSpec, %q", html.EscapeString(r.URL.Path))
}
