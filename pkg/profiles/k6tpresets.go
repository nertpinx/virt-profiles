package virtprofiles

import (
	k8sv1 "k8s.io/api/core/v1"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

func mergeDomainSpec(presetSpec *k6tv1.DomainPresetSpec, vmiSpec *k6tv1.DomainSpec) bool {
	applied := false

	if len(presetSpec.Resources.Requests) > 0 {
		if vmiSpec.Resources.Requests == nil {
			vmiSpec.Resources.Requests = k8sv1.ResourceList{}
		}

		for key, val := range presetSpec.Resources.Requests {
			current := vmiSpec.Resources.Requests[key]
			// request exceeds current setting
			if val.Cmp(current) == 1 {
				vmiSpec.Resources.Requests[key] = val
				applied = true
			}
		}
	}
	return applied
}
