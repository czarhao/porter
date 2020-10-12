package rootfs

type manager struct {
	rootPath, mntPath, writePath string
	volume                       map[string]string
}

func NewManager(rootPath, mntPath, writePath string) Manager {
	return &manager{
		rootPath:  rootPath,
		mntPath:   mntPath,
		writePath: writePath,
	}
}

func (m *manager) Create(volume, containerName, imagePath string) error {
	if err := m.createReadOnlyLayer(imagePath); err != nil {
		return err
	}
	if err := m.createWriteLayer(containerName); err != nil {
		return err
	}
	if err := m.createMountPoint(containerName, imagePath); err != nil {
		return err
	}
	return nil
}

func (m *manager) Destroy(containerName string) error {
	if err := m.deleteMountPoint(containerName); err != nil {
		return err
	}
	if err := m.deleteWriteLayer(containerName); err != nil {
		return err
	}
	return nil
}
