package virtprofiles

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

func ApplyProfiles(domSpec *libvirtxml.Domain, profiles []string) (*libvirtxml.Domain, error) {
	return domSpec, nil
}
