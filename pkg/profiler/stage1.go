/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 */

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

// ApplyPresets applies all the given presets to the stage1 domain specification
func (p *Profiler) ApplyPresets(domSpec *k6tv1.DomainSpec, presets []k6tv1.VirtualMachineInstancePreset) (*k6tv1.DomainSpec, []string, error) {
	warnings := []string{}
	ret, err := cloneDomainSpec(domSpec)
	if err != nil {
		return nil, warnings, err
	}

	domPresets, err := SortPresets(presets)
	if err != nil {
		// sorting errors are not critical for this flow
		warnings = append(warnings, fmt.Sprintf("%v", err))
	}

	err = checkPresetConflicts(domPresets)
	if err != nil {
		ret := errors.New(fmt.Sprintf("VirtualMachinePresets cannot be applied due to conflicts: %v", err))
		return nil, warnings, ret
	}

	for _, preset := range domPresets {
		applied, err := mergeDomainSpec(ret, preset.Spec.Domain)
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

func mergeDomainSpec(domSpec *k6tv1.DomainSpec, presetSpec *k6tv1.DomainPresetSpec) (bool, error) {
	presetConflicts := checkMergeConflicts(presetSpec, domSpec)
	applied := false

	if len(presetSpec.Resources.Requests) > 0 {
		if domSpec.Resources.Requests == nil {
			domSpec.Resources.Requests = k8sv1.ResourceList{}
		}
		for key, val := range presetSpec.Resources.Requests {
			curVal := domSpec.Resources.Requests[key]
			if val.Cmp(curVal) != -1 {
				domSpec.Resources.Requests[key] = val
				applied = true
			}
		}
	}
	if presetSpec.CPU != nil {
		if domSpec.CPU == nil {
			domSpec.CPU = &k6tv1.CPU{
				Model: presetSpec.CPU.Model,
			}
			applied = true
		}
		if presetSpec.CPU.Cores > domSpec.CPU.Cores {
			domSpec.CPU.Cores = presetSpec.CPU.Cores
			applied = true
		}
	}

	// TODO: handle memory
	// TODO: handle machine

	if presetSpec.Firmware != nil {
		if domSpec.Firmware == nil {
			domSpec.Firmware = &k6tv1.Firmware{}
			presetSpec.Firmware.DeepCopyInto(domSpec.Firmware)
		}
		if reflect.DeepEqual(domSpec.Firmware, presetSpec.Firmware) {
			applied = true
		}
	}
	if presetSpec.Clock != nil {
		if domSpec.Clock == nil {
			domSpec.Clock = &k6tv1.Clock{}
			domSpec.Clock.ClockOffset = presetSpec.Clock.ClockOffset
		}
		if reflect.DeepEqual(domSpec.Clock, presetSpec.Clock) {
			applied = true
		}

		if presetSpec.Clock.Timer != nil {
			if domSpec.Clock.Timer == nil {
				domSpec.Clock.Timer = &k6tv1.Timer{}
				presetSpec.Clock.Timer.DeepCopyInto(domSpec.Clock.Timer)
			}
			if reflect.DeepEqual(domSpec.Clock.Timer, presetSpec.Clock.Timer) {
				applied = true
			}
		}
	}
	if presetSpec.Features != nil {
		if domSpec.Features == nil {
			domSpec.Features = &k6tv1.Features{}
			presetSpec.Features.DeepCopyInto(domSpec.Features)
		}
		if reflect.DeepEqual(domSpec.Features, presetSpec.Features) {
			applied = true
		}
	}
	if presetSpec.Devices.Watchdog != nil {
		if domSpec.Devices.Watchdog == nil {
			domSpec.Devices.Watchdog = &k6tv1.Watchdog{}
			presetSpec.Devices.Watchdog.DeepCopyInto(domSpec.Devices.Watchdog)
		}
		if reflect.DeepEqual(domSpec.Devices.Watchdog, presetSpec.Devices.Watchdog) {
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

func checkMergeConflicts(presetSpec *k6tv1.DomainPresetSpec, vmiSpec *k6tv1.DomainSpec) error {
	errors := []error{}

	// resource request never conflicts: we pick the union of the requests, and the larger value among overlapping requests

	// same for cpu cores.

	if presetSpec.CPU != nil && vmiSpec.CPU != nil {
		if !reflect.DeepEqual(presetSpec.CPU.Model, vmiSpec.CPU.Model) {
			errors = append(errors, fmt.Errorf("spec.cpu.model: %v != %v", presetSpec.CPU.Model, vmiSpec.CPU.Model))
		}
	}

	// TODO: check memory

	// TODO: check machine

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

func cloneDomainSpec(dom *k6tv1.DomainSpec) (*k6tv1.DomainSpec, error) {
	ret := &k6tv1.DomainSpec{}
	data, _ := json.Marshal(dom)
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
