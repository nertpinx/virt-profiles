// Package virtprofiles implements the VM profile processing
//
// A virt profile is a set of changes to be performed to a given representation of a VM.
// For the purposes of this package, a representation of a VM is either a kubevirt API
// DomainSpec object, or a libvirt XML document.
package virtprofiles

import (
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

// Catalogue manages a collection of virt profiles.
type Catalogue struct {
}

func NewCatalogue(profilesDir string) (*Catalogue, error) {
	// TODO: make sure profilesDir is abspath
	return &Catalogue{}, nil
}

// Names return the names of all the profiles in the Catalogue
// the profile names are treated as opaque strings (e.g. they
// don't have an implicit meaning) that can be used later to
// refer to profiles.
func (c *Catalogue) Names() ([]string, error) {
	entries := []string{}
	return entries, nil
}

func (c *Catalogue) AddPreset(preset k6tv1.DomainPresetSpec) error {
	return nil
}

func (c *Catalogue) Get(name string) (interface{}, error) {
	return nil, nil
}

func (c *Catalogue) GetAll(names []string) ([]interface{}, error) {
	return nil, nil
}
