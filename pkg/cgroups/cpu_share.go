package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"porter/pkg/utils"
)

type CpuShareSubsystem struct{}

var _ Subsystem = &CpuShareSubsystem{}

func (cpu *CpuShareSubsystem) Name() string {
	return "cpu"
}

func (cpu *CpuShareSubsystem) Set(cgroupPath string, resConfig *utils.Limit) error {
	subsystemPath, err := getCgroupPath(cpu.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if resConfig.CpuShare != "" {
		if err := ioutil.WriteFile(path.Join(subsystemPath, "cpu.shares"), []byte(resConfig.CpuShare), 0644); err != nil {
			return fmt.Errorf("set subsystem cpu share fail %v", err)
		}
	}
	return nil
}

func (cpu *CpuShareSubsystem) Add(cgroupPath string, pid int) error {
	subsystemPath, err := getCgroupPath(cpu.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("get subsystem %s error: %v", cgroupPath, err)
	}
	if err := writeCgroupProc(subsystemPath, pid); err != nil {
		return fmt.Errorf("set cpu_share subsystem proc fail %v", err)
	}
	return nil
}

func (cpu *CpuShareSubsystem) Remove(cgroupPath string) error {
	subsystemPath, err := getCgroupPath(cpu.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("remove subsystem %s fail %v", cpu.Name(), err)
	}
	return os.RemoveAll(subsystemPath)
}
