package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"porter/pkg/utils"
)

type MemorySubsystem struct{}

var _ Subsystem = &MemorySubsystem{}

func (mem *MemorySubsystem) Name() string {
	return "memory"
}

func (mem *MemorySubsystem) Set(cgroupPath string, resConfig *utils.Limit) error {
	subsystemPath, err := getCgroupPath(mem.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if resConfig.MemoryLimit != "" {
		if err := ioutil.WriteFile(path.Join(subsystemPath, "memory.limit_in_bytes"), []byte(resConfig.MemoryLimit), 0644); err != nil {
			return fmt.Errorf("set subsystem memory fail %v", err)
		}
	}
	return nil
}

func (mem *MemorySubsystem) Add(cgroupPath string, pid int) error {
	subsystemPath, err := getCgroupPath(mem.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("get subsystem %s error: %v", cgroupPath, err)
	}
	if err := writeCgroupProc(subsystemPath, pid); err != nil {
		return fmt.Errorf("set memory subsystem proc fail %v", err)
	}
	return nil
}

func (mem *MemorySubsystem) Remove(cgroupPath string) error {
	subsystemPath, err := getCgroupPath(mem.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("remove subsystem %s fail %v", mem.Name(), err)
	}
	return os.RemoveAll(subsystemPath)
}
