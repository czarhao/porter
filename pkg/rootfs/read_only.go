package rootfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"porter/pkg/utils"
)

func (m *manager) createReadOnlyLayer(imagePath string) error {
	imageName := utils.Path2name(imagePath, true)
	unTarFolderUrl := path.Join(m.rootPath, imageName)
	exist, err := utils.PathExists(unTarFolderUrl)
	if err != nil {
		return fmt.Errorf("fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			return fmt.Errorf("mkdir %s error %v", unTarFolderUrl, err)
		}
		if _, err := exec.Command("tar", "-xvf", imagePath, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			return fmt.Errorf("untar dir %s error %v", unTarFolderUrl, err)
		}
	}
	return nil
}
