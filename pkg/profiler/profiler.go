package virtprofiles

import (
	k8sv1 "k8s.io/api/core/v1"
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

type Profiler struct {
	secrets        map[string]*k8sv1.Secret
	virtualMachine *k6tv1.VirtualMachineInstance
	baseDiskPath   string
}

func (p *Profiler) AddSecret(key string, value *k8sv1.Secret) *Profiler {
	p.secrets[key] = value
	return p
}

func (p *Profiler) SetVirtualMachine(vm *k6tv1.VirtualMachineInstance) *Profiler {
	p.virtualMachine = vm
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
		secrets:      make(map[string]*k8sv1.Secret),
		baseDiskPath: basePath,
	}
}
