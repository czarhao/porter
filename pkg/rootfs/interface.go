package rootfs

type Manager interface {
	Create(volume, containerName, imagePath string) error
	Destroy(containerName string) error
}
