package actions

import (
	"porter/pkg/container"
	"porter/pkg/utils"

	"github.com/urfave/cli/v2"
)

var Run = func(c *cli.Context) error {
	path := c.Args().First()
	cfg, err := utils.Generate(path)
	if err != nil {
		utils.Logger.Fatalf("parsing %s fail, %v", path, err)
	}

	var (
		containerCfg = cfg.Container
		imageCfg     = cfg.Image

		limitCfg     = cfg.Limit
	)

	proc, err := container.NewManager(containerCfg, imageCfg)
	if err != nil {
		utils.Logger.Fatalf("failed to create container, %v", err)
	}

	if err := proc.SetLimit(limitCfg); err != nil {
		utils.Logger.Fatalf("failed to set cgroup or start, %v", err)
	} else {
		utils.Logger.Infof("set cgroup successfully")
	}

	defer func() {
		if err := proc.DestroyCgroup(); err != nil {
			utils.Logger.Fatalf("failed to clear cgroup, %v", err)
		} else {
			utils.Logger.Infof("clean cgroup successfully")
		}
	}()

	if err := proc.RunCmd([]string{imageCfg.Run}); err != nil {
		utils.Logger.Fatalf("send user command to container proc error, %v", err)
	}

	if containerCfg.Tty {
		_ = proc.WaitProcEnd()
	}

	if err := proc.DestroyMount(); err != nil {
		utils.Logger.Fatalf("failed to clear mount, %v", err)
	} else {
		utils.Logger.Infof("clear mount successfully")
	}

	return nil
}
