package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"porter/pkg/utils"
	"strings"
	"syscall"

	"porter/pkg/cgroups"
	"porter/pkg/rootfs"
)

type manager struct {
	cmd   *exec.Cmd
	wPipe *os.File

	cgroup cgroups.Manager
	rootfs rootfs.Manager

	container *utils.Container

	tty bool
}

func NewManager(container *utils.Container, image *utils.Image) (Manager, error) {
	rPipe, wPipe, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe, %v", err)
	}

	proc := exec.Command("/proc/self/exe", "init")

	proc.SysProcAttr = &syscall.SysProcAttr{Cloneflags: CLONE_FLAGS}

	if container.Tty {
		proc.Stdin = os.Stdin
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
	}

	proc.ExtraFiles = []*os.File{rPipe}

	rootfsManager := rootfs.NewManager(
		image.Layer.RootURL,
		image.Layer.MntURL,
		image.Layer.WriterUrl,
	)

	if err := rootfsManager.Create(
		image.Layer.Volume,
		container.Name,
		image.Path,
	); err != nil {
		return nil, err
	}

	proc.Dir = path.Join(image.Layer.MntURL, container.Name)

	return &manager{
		cmd:           proc,
		wPipe:         wPipe,
		container: container,
		cgroup: cgroups.NewManager(container.Name),
		rootfs: rootfsManager,
	}, nil
}

func (p *manager) SetLimit(limit *utils.Limit) error {
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("start parent proccess fail, %v", err)
	}
	if err := p.cgroup.Set(limit); err != nil {
		return err
	}
	return p.cgroup.Add(p.cmd.Process.Pid)
}

func (p *manager) RunCmd(cmds []string) error {
	command := strings.Join(cmds, " ")
	utils.Logger.Infof("Porter run command is %s", command)
	if _, err := p.wPipe.WriteString(command); err != nil {
		return err
	}
	return p.wPipe.Close()}

func (p *manager) WaitProcEnd() error {
	return p.cmd.Wait()
}

func (p *manager) DestroyMount() error {
	return p.rootfs.Destroy(p.container.Name)
}

func (p *manager) DestroyCgroup() error {
	return p.cgroup.Destroy()
}
