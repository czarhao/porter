package mount

import (
	"fmt"
	"os"
	"syscall"
)

const (
	procMountFlags  = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	tmpfsMountFlags = syscall.MS_NOSUID | syscall.MS_STRICTATIME
)

func Set() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get pwd fail, error: %v", err)
	}

	if err := pivotRoot(pwd); err != nil {
		return err
	}

	//mount proc
	if err := syscall.Mount("proc", "/proc", "proc",
		uintptr(procMountFlags), ""); err != nil {
		return fmt.Errorf("mount /proc fail, error: %v", err)
	}

	//mount tmpfs
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs",
		uintptr(tmpfsMountFlags), "mode=755"); err != nil {
		return fmt.Errorf("mount /dev fail, error: %v", err)
	}
	return nil
}
