package virtprofiles

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

func TranslateSpecs(domSpec *k6tv1.DomainSpec) (*libvirtxml.Domain, error) {
	ret := &libvirtxml.Domain{}
	return ret, nil
}
