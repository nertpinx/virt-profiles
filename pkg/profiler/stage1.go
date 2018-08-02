package virtprofiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	k8sv1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

func ApplyPresets(domSpec *k6tv1.DomainSpec, presets []k6tv1.VirtualMachineInstancePreset) (*k6tv1.DomainSpec, []string, error) {
	warnings := []string{}
	ret, err := cloneDomainSpec(domSpec)
	if err != nil {
		return nil, warnings, err
	}

	err = checkPresetConflicts(presets)
	if err != nil {
		ret := errors.New(fmt.Sprintf("VirtualMachinePresets cannot be applied due to conflicts: %v", err))
		return nil, warnings, ret
	}

	for _, preset := range presets {
		applied, err := mergeDomainSpec(preset.Spec.Domain, ret)
		if err != nil {
			msg := fmt.Sprintf("Unable to apply VirtualMachineInstancePreset '%s': %v", preset.Name, err)
			if applied {
				msg = fmt.Sprintf("Some settings were not applied for VirtualMachineInstancePreset '%s': %v", preset.Name, err)
			}

			warnings = append(warnings, msg)
		}
	}
	return ret, warnings, nil
}

func checkMergeConflicts(presetSpec *k6tv1.DomainPresetSpec, vmiSpec *k6tv1.DomainSpec) error {
	errors := []error{}
	if len(presetSpec.Resources.Requests) > 0 {
		for key, presetReq := range presetSpec.Resources.Requests {
			if vmiReq, ok := vmiSpec.Resources.Requests[key]; ok {
				if presetReq != vmiReq {
					errors = append(errors, fmt.Errorf("spec.resources.requests[%s]: %v != %v", key, presetReq, vmiReq))
				}
			}
		}
	}
	if presetSpec.CPU != nil && vmiSpec.CPU != nil {
		if !reflect.DeepEqual(presetSpec.CPU, vmiSpec.CPU) {
			errors = append(errors, fmt.Errorf("spec.cpu: %v != %v", presetSpec.CPU, vmiSpec.CPU))
		}
	}
	if presetSpec.Firmware != nil && vmiSpec.Firmware != nil {
		if !reflect.DeepEqual(presetSpec.Firmware, vmiSpec.Firmware) {
			errors = append(errors, fmt.Errorf("spec.firmware: %v != %v", presetSpec.Firmware, vmiSpec.Firmware))
		}
	}
	if presetSpec.Clock != nil && vmiSpec.Clock != nil {
		if !reflect.DeepEqual(presetSpec.Clock.ClockOffset, vmiSpec.Clock.ClockOffset) {
			errors = append(errors, fmt.Errorf("spec.clock.clockoffset: %v != %v", presetSpec.Clock.ClockOffset, vmiSpec.Clock.ClockOffset))
		}
		if presetSpec.Clock.Timer != nil && vmiSpec.Clock.Timer != nil {
			if !reflect.DeepEqual(presetSpec.Clock.Timer, vmiSpec.Clock.Timer) {
				errors = append(errors, fmt.Errorf("spec.clock.timer: %v != %v", presetSpec.Clock.Timer, vmiSpec.Clock.Timer))
			}
		}
	}
	if presetSpec.Features != nil && vmiSpec.Features != nil {
		if !reflect.DeepEqual(presetSpec.Features, vmiSpec.Features) {
			errors = append(errors, fmt.Errorf("spec.features: %v != %v", presetSpec.Features, vmiSpec.Features))
		}
	}
	if presetSpec.Devices.Watchdog != nil && vmiSpec.Devices.Watchdog != nil {
		if !reflect.DeepEqual(presetSpec.Devices.Watchdog, vmiSpec.Devices.Watchdog) {
			errors = append(errors, fmt.Errorf("spec.devices.watchdog: %v != %v", presetSpec.Devices.Watchdog, vmiSpec.Devices.Watchdog))
		}
	}

	if len(errors) > 0 {
		return utilerrors.NewAggregate(errors)
	}
	return nil
}

func mergeDomainSpec(presetSpec *k6tv1.DomainPresetSpec, vmiSpec *k6tv1.DomainSpec) (bool, error) {
	presetConflicts := checkMergeConflicts(presetSpec, vmiSpec)
	applied := false

	if len(presetSpec.Resources.Requests) > 0 {
		if vmiSpec.Resources.Requests == nil {
			vmiSpec.Resources.Requests = k8sv1.ResourceList{}
			for key, val := range presetSpec.Resources.Requests {
				vmiSpec.Resources.Requests[key] = val
			}
		}
		if reflect.DeepEqual(vmiSpec.Resources.Requests, presetSpec.Resources.Requests) {
			applied = true
		}
	}
	if presetSpec.CPU != nil {
		if vmiSpec.CPU == nil {
			vmiSpec.CPU = &k6tv1.CPU{}
			presetSpec.CPU.DeepCopyInto(vmiSpec.CPU)
		}
		if reflect.DeepEqual(vmiSpec.CPU, presetSpec.CPU) {
			applied = true
		}
	}
	if presetSpec.Firmware != nil {
		if vmiSpec.Firmware == nil {
			vmiSpec.Firmware = &k6tv1.Firmware{}
			presetSpec.Firmware.DeepCopyInto(vmiSpec.Firmware)
		}
		if reflect.DeepEqual(vmiSpec.Firmware, presetSpec.Firmware) {
			applied = true
		}
	}
	if presetSpec.Clock != nil {
		if vmiSpec.Clock == nil {
			vmiSpec.Clock = &k6tv1.Clock{}
			vmiSpec.Clock.ClockOffset = presetSpec.Clock.ClockOffset
		}
		if reflect.DeepEqual(vmiSpec.Clock, presetSpec.Clock) {
			applied = true
		}

		if presetSpec.Clock.Timer != nil {
			if vmiSpec.Clock.Timer == nil {
				vmiSpec.Clock.Timer = &k6tv1.Timer{}
				presetSpec.Clock.Timer.DeepCopyInto(vmiSpec.Clock.Timer)
			}
			if reflect.DeepEqual(vmiSpec.Clock.Timer, presetSpec.Clock.Timer) {
				applied = true
			}
		}
	}
	if presetSpec.Features != nil {
		if vmiSpec.Features == nil {
			vmiSpec.Features = &k6tv1.Features{}
			presetSpec.Features.DeepCopyInto(vmiSpec.Features)
		}
		if reflect.DeepEqual(vmiSpec.Features, presetSpec.Features) {
			applied = true
		}
	}
	if presetSpec.Devices.Watchdog != nil {
		if vmiSpec.Devices.Watchdog == nil {
			vmiSpec.Devices.Watchdog = &k6tv1.Watchdog{}
			presetSpec.Devices.Watchdog.DeepCopyInto(vmiSpec.Devices.Watchdog)
		}
		if reflect.DeepEqual(vmiSpec.Devices.Watchdog, presetSpec.Devices.Watchdog) {
			applied = true
		}
	}
	return applied, presetConflicts
}

func cloneDomainSpec(dom *k6tv1.DomainSpec) (*k6tv1.DomainSpec, error) {
	ret := &k6tv1.DomainSpec{}
	data, _ := json.Marshal(dom)
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Compare the domain of every preset to ensure they can all be applied cleanly
func checkPresetConflicts(presets []k6tv1.VirtualMachineInstancePreset) error {
	errors := []error{}
	visitedPresets := []k6tv1.VirtualMachineInstancePreset{}
	for _, preset := range presets {
		for _, visited := range visitedPresets {
			visitedDomain := &k6tv1.DomainSpec{}
			domainByte, _ := json.Marshal(visited.Spec.Domain)
			err := json.Unmarshal(domainByte, &visitedDomain)
			if err != nil {
				return err
			}

			err = checkMergeConflicts(preset.Spec.Domain, visitedDomain)
			if err != nil {
				errors = append(errors, fmt.Errorf("presets '%s' and '%s' conflict: %v", preset.Name, visited.Name, err))
			}
		}
		visitedPresets = append(visitedPresets, preset)
	}
	if len(errors) > 0 {
		return utilerrors.NewAggregate(errors)
	}
	return nil
}
