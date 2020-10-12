package rootfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"porter/pkg/utils"
)

func (m *manager) createMountPoint(containerName, imagePath string) error {
	mntDir := path.Join(m.mntPath, containerName)
	if err := os.MkdirAll(mntDir, 0777); err != nil {
		return fmt.Errorf("mkdir %s error %v", mntDir, err)
	}

	var (
		imageName = utils.Path2name(imagePath, true)
		writeLayer    = path.Join(m.writePath, containerName)
		imageLocation = path.Join(m.rootPath, imageName)
		mntURL        = path.Join(m.mntPath, containerName)
	)

	dirs := "dirs=" + writeLayer + ":" + imageLocation
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL).CombinedOutput()
	if err != nil {
		return fmt.Errorf("run command for creating mount point failed %v", err)
	}
	return nil
}

func (m *manager) deleteMountPoint(containerName string) error {
	mntURL := path.Join(m.mntPath, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unmount %s error %v", mntURL, err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		return fmt.Errorf("remove mountpoint dir %s error %v", mntURL, err)
	}
	return nil
}
