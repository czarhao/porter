package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

// 下面的内容参考了 docker 的 pivotRoot 函数实现...
func pivotRoot(path string) (err error) {
	// if the engine is running in a user namespace we need to use actual chroot
	if isRunningInUserNS() {
		return chroot(path)
	}

	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		return fmt.Errorf("error creating mount namespace before pivot: %v", err)
	}

	if err := makeRSlave("/"); err != nil {
		return err
	}

	if mounted, _ := isMounted(path); !mounted {
		if err := mount(path, path, "bind", "rbind,rw"); err != nil {
			return chroot(path)
		}
	}

	// setup oldRoot for pivot_root
	pivotDir, err := ioutil.TempDir(path, ".pivot_root")
	if err != nil {
		return fmt.Errorf("error setting up pivot dir: %v", err)
	}

	var mounted bool
	defer func() {
		if mounted {
			// make sure pivotDir is not mounted before we try to remove it
			if errCleanup := syscall.Unmount(pivotDir, syscall.MNT_DETACH); errCleanup != nil {
				if err == nil {
					err = errCleanup
				}
				return
			}
		}

		errCleanup := os.Remove(pivotDir)
		// pivotDir doesn't exist if pivot_root failed and chroot+chdir was successful
		// because we already cleaned it up on failed pivot_root
		if errCleanup != nil && !os.IsNotExist(errCleanup) {
			errCleanup = fmt.Errorf("error cleaning up after pivot: %v", errCleanup)
			if err == nil {
				err = errCleanup
			}
		}
	}()

	if err := syscall.PivotRoot(path, pivotDir); err != nil {
		// If pivot fails, fall back to the normal chroot after cleaning up temp dir
		if err := os.Remove(pivotDir); err != nil {
			return fmt.Errorf("error cleaning up after failed pivot: %v", err)
		}
		return chroot(path)
	}
	mounted = true

	// This is the new path for where the old root (prior to the pivot) has been moved to
	// This dir contains the layer of the caller, which we need to remove so it is not visible during extraction
	pivotDir = filepath.Join("/", filepath.Base(pivotDir))

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("error changing to new root: %v", err)
	}

	// Make the pivotDir (where the old root lives) private so it can be unmounted without propagating to the host
	if err := syscall.Mount("", pivotDir, "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("error making old root private after pivot: %v", err)
	}

	// Now unmount the old root so it's no longer visible from the new root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("error while unmounting old root after pivot: %v", err)
	}
	mounted = false
	return nil
}
