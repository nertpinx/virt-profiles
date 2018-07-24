package virtprofiles

import "io/ioutil"

type Catalogue struct {
	profilesDir string
}

func NewCatalogue(profilesDir string) (*Catalogue, error) {
	// TODO: make sure profilesDir is abspath
	return &Catalogue{profilesDir: profilesDir}, nil
}

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
