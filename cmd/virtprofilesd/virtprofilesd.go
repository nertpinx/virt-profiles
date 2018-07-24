package main

import (
	"fmt"
	"log"
	"net/http"

	flag "github.com/spf13/pflag"
)

func main() {
	conf := Config{}
	conf.ParseFlags()

	log.Printf("profiles from %s", conf.Profiles)
	app, err := NewProfilesApp(conf.Profiles)
	if err != nil {
		log.Fatal("%v", err)
	}

	log.Printf("listening on %s", conf.ListenAddress())
	log.Fatal(http.ListenAndServe(conf.ListenAddress(), app))
}

type Config struct {
	Host     string
	Port     int
	Profiles string
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.Host, "host", "localhost", "set the interface to listen to")
	flag.IntVar(&c.Port, "port", 8080, "set the port to listen to")
	flag.StringVar(&c.Profiles, "profiles", "/usr/share/virt-profiles", "set the libvirt profiles directory")
	flag.Parse()
}

func (c *Config) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
