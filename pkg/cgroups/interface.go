package cgroups

import "porter/pkg/utils"

// subsystem
type Subsystem interface {
	Name() string
	Set(cgroup string, resConfig *utils.Limit) error
	Add(cgroup string, pid int) error
	Remove(cgroup string) error
}
