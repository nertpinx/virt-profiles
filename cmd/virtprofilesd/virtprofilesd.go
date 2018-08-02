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

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fromanirh/virt-profiles/internal/pkg/profilerapp"
	flag "github.com/spf13/pflag"
)

func main() {
	conf := Config{}
	conf.ParseFlags()

	log.Printf("profiles from %s", conf.Profiles)
	app, err := profilerapp.NewProfilerApp(conf.Profiles)
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
