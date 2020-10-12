package actions

import (
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"porter/pkg/mount"
	"porter/pkg/utils"
	"syscall"
)

var Init = func(c *cli.Context) error {
	cmds, err := utils.ReadUserCommand()

	if err != nil {
		utils.Logger.Fatalf("run container get user command error, %v", err)
	} else if len(cmds) == 0 {
		utils.Logger.Fatalf("run container get user command error, len(cmd) == 0")
	}

	if err := mount.Set(); err != nil {
		utils.Logger.Fatalf("mount fail, %v", err)
	}

	path, err := exec.LookPath(cmds[0])
	if err != nil {
		utils.Logger.Fatalf("look path fail, %v", err)
	}

	if err := syscall.Exec(path, cmds[0:], os.Environ()); err != nil {
		utils.Logger.Fatalf("exec container command fail", err)
	}

	return nil
}
