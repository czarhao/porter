package cgroups

import "porter/pkg/utils"

type Manager interface {
	Add(pid int) error
	Set(resource *utils.Limit) error
	Destroy() error
}

type subsystemManager struct {
	path       string
	subsystems []Subsystem
}

func NewManager(path string) *subsystemManager {
	return &subsystemManager{
		path: path,
		subsystems: []Subsystem{
			&MemorySubsystem{},
			&CpuShareSubsystem{},
		},
	}
}

func (cm *subsystemManager) Add(pid int) error {
	for _, v := range cm.subsystems {
		if err := v.Add(cm.path, pid); err != nil {
			return err
		}
	}
	return nil
}

func (cm *subsystemManager) Set(resource *utils.Limit) error {
	for _, v := range cm.subsystems {
		if err := v.Set(cm.path, resource); err != nil {
			return err
		}
	}
	return nil
}

func (cm *subsystemManager) Destroy() error {
	for _, v := range cm.subsystems {
		if err := v.Remove(cm.path); err != nil {
			return err
		}
	}
	return nil
}
