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

package virtprofiles

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

// ApplyProfiles applies all the given XML profiles to the stage3 domain specification
func (p *Profiler) ApplyProfiles(domSpec *libvirtxml.Domain, profiles []string) (*libvirtxml.Domain, error) {
	// Does'nt really apply them until we figure out how to implement XML profiles.
	return domSpec, nil
}

// Complete fills the unspecified backend settings with optimal values
func (p *Profiler) Complete(domSpec *libvirtxml.Domain) (*libvirtxml.Domain, error) {
	return domSpec, nil
}
