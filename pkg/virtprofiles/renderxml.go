package virtprofiles

import (
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

// RenderXML applies all the given profiles to the given domain spec,
// and returns a libvirt domain XML document.
func (c *Catalogue) RenderXML(domain k6tv1.DomainSpec, profiles []string) (string, error) {
	return "", nil
}
