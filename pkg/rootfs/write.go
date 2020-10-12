package rootfs

import (
	"fmt"
	"os"
	"path"
)

func (m *manager) createWriteLayer(containerName string) error {
	writeURL := path.Join(m.writePath, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		return fmt.Errorf("mkdir write layer dir %s error. %v", writeURL, err)
	}
	return nil
}

func (m *manager) deleteWriteLayer(containerName string) error {
	writeURL := path.Join(m.writePath, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		return fmt.Errorf("remove writeLayer dir %s error %v", writeURL, err)
	}
	return nil
}
