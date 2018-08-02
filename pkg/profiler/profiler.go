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
	k8sv1 "k8s.io/api/core/v1"
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

const priorityMarking = "virtualmachineinstancepresets.admission.kubevirt.io/priority"

type Profiler struct {
	secrets           map[string]*k8sv1.Secret
	virtualMachine    *k6tv1.VirtualMachineInstance
	baseDiskPath      string
	sortingAnnotation string
}

func (p *Profiler) AddSecret(key string, value *k8sv1.Secret) *Profiler {
	p.secrets[key] = value
	return p
}

func (p *Profiler) SetVirtualMachine(vm *k6tv1.VirtualMachineInstance) *Profiler {
	p.virtualMachine = vm
	return p
}

func (p *Profiler) SetPriorityMarking(marking string) *Profiler {
	p.sortingAnnotation = priorityMarking
	return p
}

func (p *Profiler) SetBaseDiskPath(path string) *Profiler {
	p.baseDiskPath = path
	return p
}

func (p *Profiler) BaseDiskPath() string {
	return p.baseDiskPath
}

func NewProfiler(basePath string) *Profiler {
	return &Profiler{
		secrets:           make(map[string]*k8sv1.Secret),
		baseDiskPath:      basePath,
		sortingAnnotation: priorityMarking,
	}
}
