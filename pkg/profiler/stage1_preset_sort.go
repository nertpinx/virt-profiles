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
	"errors"
	"sort"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

// sortPresets sorts and returns a slice of VirtualMachinePresets, using optional annotations.
func (p *Profiler) SortPresets(presets []k6tv1.VirtualMachineInstancePreset) ([]k6tv1.VirtualMachineInstancePreset, error) {
	err := checkAnnotations(presets)
	if err != nil {
		return presets, err
	}
	sort.Stable(byPriority{Presets: presets, Annotation: p.Annotation})
	return presets, nil
}

func checkAnnotations(presets []k6tv1.VirtualMachineInstancePreset) error {
	for _, preset := range presets {
		if preset.Annotations == nil {
			return errors.New("preset %v lacks annotations", preset.Name)
		}
		_, ok := preset.Annotations[p.Annotation]
		if !ok {
			return errors.New("preset %v lacks priority annotation", preset.Name)
		}
	}
	return nil
}

type byPriority struct {
	Presets    []k6tv1.VirtualMachineInstancePreset
	Annotation string
}

// sort.Interface.
func (p *byPriority) Len() int {
	return len(p.Presets)
}
func (p *byPriority) Swap(i, j int) {
	p.Presets[i], p.Presets[j] = p.Presets[j], p.Presets[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (p *byPriority) Less(i, j int) bool {
	var prio1, prio2 int
	var err1, err2 error
	if value1, ok := p.Presets[i].Annotations[p.Annotation]; ok {
		prio1, err1 = strconv.Atoi(value1)
	}
	if value2, ok := p.Presets[j].Annotations[p.Annotation]; ok {
		prio2, err2 = strconv.Atoi(value2)
	}
	// only if we succesfully parsed both priorities we can make a meaningful comparation
	if err1 != nil || err2 != nil {
		return true
	}
	// intentionally using ">" here. The higher the priority, the sooner the preset should
	// be in the sequence, so the earlier will be applied
	return prio1 > prio2
}
