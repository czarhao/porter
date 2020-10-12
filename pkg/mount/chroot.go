package mount

import (
	"fmt"
	"syscall"

	"porter/pkg/utils"
)

func chroot(path string) error {
	utils.Logger.Errorf("use pivot_root fail, try chroot")
	if err := syscall.Chroot(path); err != nil {
		return fmt.Errorf("error after fallback to chroot: %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("error changing to new root after chroot: %v", err)
	}
	return nil
}
