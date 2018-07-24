package virtprofiles

import (
	k8sv1 "k8s.io/api/core/v1"
	//	k8sres "k8s.io/apimachinery/pkg/api/resource"
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

// ComputeResources applies all the given profiles to the given domain spec,
// and returns all the resources needed by the resulting domain.
func (c *Catalogue) ComputeResources(domain k6tv1.DomainSpec, profiles []string) (k8sv1.ResourceList, error) {
	res := make(k8sv1.ResourceList)
	return res, nil
}
