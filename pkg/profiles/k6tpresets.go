package virtprofiles

import (
	"encoding/json"
	"fmt"
	"reflect"

	k8sv1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

func ApplyDomainSpecPresets(domSpec k6tv1.DomainSpec, presets []k6tv1.VirtualMachineInstancePreset) (k6tv1.DomainSpec, error) {
	errors := []error{}
	resultSpec, err := cloneDomainSpec(domSpec)
	if err != nil {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return resultSpec, utilerrors.NewAggregate(errors)
	}
	return resultSpec, nil
}

func cloneDomainSpec(domSpec k6tv1.DomainSpec) (k6tv1.DomainSpec, error) {
	clonedSpec := k6tv1.DomainSpec{}
	domainBytes, _ := json.Marshal(domSpec)
	err := json.Unmarshal(domainBytes, &clonedSpec)
	return clonedSpec, err
}

func checkMergeConflicts(presetSpec *k6tv1.DomainPresetSpec, vmiSpec *k6tv1.DomainSpec) error {
	errors := []error{}

	// resources never conflict: we pick the union of all the requests

	if presetSpec.CPU != nil && vmiSpec.CPU != nil {
		// same with CPU cores, we pick the biggest requirement
		// TODO: compare model and see if they are compatible. Review if libvirt APIs may help.
		if !reflect.DeepEqual(presetSpec.CPU.Model, vmiSpec.CPU.Model) {
			errors = append(errors, fmt.Errorf("spec.cpu.model: %v != %v", presetSpec.CPU.Model, vmiSpec.CPU.Model))
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
	// TODO: we can be more specific here. ACPI and APIC are flags, so we can relax the constraint here;
	// OTOH timers and HyperV are (practically) mutually exclusive.
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

	// TODO: check Memory, Machine fields?

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
	if presetSpec.CPU != nil {
		if vmiSpec.CPU == nil {
			vmiSpec.CPU = &k6tv1.CPU{}
			vmiSpec.CPU.Model = presetSpec.CPU.Model
			applied = true
		}

		if presetSpec.CPU.Cores > vmiSpec.CPU.Cores {
			vmiSpec.CPU.Cores = presetSpec.CPU.Cores
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

/*
func applyPresets(vmi *k6tv1.VirtualMachineInstance, presets []k6tv1.VirtualMachineInstancePreset, recorder record.EventRecorder) bool {
	logger := log.Log
	err := checkPresetConflicts(presets)
	if err != nil {
		msg := fmt.Sprintf("VirtualMachinePresets cannot be applied due to conflicts: %v", err)
		recorder.Event(vmi, k8sv1.EventTypeWarning, k6tv1.PresetFailed.String(), msg)
		logger.Object(vmi).Error(msg)
		return false
	}

	for _, preset := range presets {
		applied, err := mergeDomainSpec(preset.Spec.Domain, &vmi.Spec.Domain)
		if err != nil {
			msg := fmt.Sprintf("Unable to apply VirtualMachineInstancePreset '%s': %v", preset.Name, err)
			if applied {
				msg = fmt.Sprintf("Some settings were not applied for VirtualMachineInstancePreset '%s': %v", preset.Name, err)
			}

			recorder.Event(vmi, k8sv1.EventTypeNormal, k6tv1.Override.String(), msg)
			logger.Object(vmi).Info(msg)
		}
		if applied {
			annotateVMI(vmi, preset)
		}
	}
	return true
}
*/
