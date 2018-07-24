package virtprofiles

type Catalogue struct {
	profilesDir string
}

func NewCatalogue(profilesDir string) (*Catalogue, error) {
	return &Catalogue{profilesDir: profilesDir}, nil
}

func (c *Catalogue) Names() []string {
	keys := []string{}
	return keys
}
