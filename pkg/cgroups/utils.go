package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	cgroupProcesses = "cgroup.procs"
)

func getCgroupPath(subsystem, cgroupPath string, create bool) (string, error) {
	cgroupRoot, err := findCgroupMountPoint(subsystem)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path.Join(cgroupRoot, cgroupPath))
	if create && os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
			return "", fmt.Errorf("subsystem path error %v", err)
		}
	}
	if err == nil {
		return path.Join(cgroupRoot, cgroupPath), nil
	}
	return "", fmt.Errorf("subsystem path error %v", err)
}

func findCgroupMountPoint(subsystem string) (string, error) {
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Split(text, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4], nil
			}
		}
	}
	return "", scanner.Err()
}

func writeCgroupProc(dir string, pid int) error {
	if dir == "" {
		return fmt.Errorf("no such directory for %s", cgroupProcesses)
	}
	if pid == -1 {
		return nil
	}
	cgroupProcessesFile, err := os.OpenFile(filepath.Join(dir, cgroupProcesses), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return fmt.Errorf("failed to write %v to %v: %v", pid, cgroupProcesses, err)
	}
	defer cgroupProcessesFile.Close()
	_, err = cgroupProcessesFile.WriteString(strconv.Itoa(pid))
	return err
}
