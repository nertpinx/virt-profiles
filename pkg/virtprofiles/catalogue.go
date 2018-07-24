package virtprofiles

import "io/ioutil"

// Catalogue manages a collection of virt profiles.
type Catalogue struct {
	profilesDir string
}

func NewCatalogue(profilesDir string) (*Catalogue, error) {
	// TODO: make sure profilesDir is abspath
	return &Catalogue{profilesDir: profilesDir}, nil
}

// Names return the names of all the profiles in the Catalogue
// the profile names are treated as opaque strings (e.g. they
// don't have an implicit meaning) that can be used later to
// refer to profiles.
func (c *Catalogue) Names() ([]string, error) {
	entries := []string{}
	files, err := ioutil.ReadDir(c.profilesDir)
	if err != nil {
		return entries, err
	}
	for _, file := range files {
		entries = append(entries, file.Name())
	}
	return entries, nil
}
