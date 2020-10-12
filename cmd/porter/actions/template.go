package actions

import (
	"github.com/urfave/cli/v2"
	"os"
	"porter/pkg/utils"
)

const YAML_TEMPLATE = `
container:
  name: dev_container
  tty: true
image:
  path: /home/czarhao/Code/go/porter/image/busybox.tar
  run: /bin/sh
  layer:
    root_url: /home/czarhao/tmp/root
    mnt_url: /home/czarhao/tmp/mnt
    writer_url: /home/czarhao/tmp/writer
limit:
  memory: 50m
  cpu_share: 512
`

var Template = func(c *cli.Context) error {
	filename := c.Args().First()
	file, err := os.Create(filename)
	if err != nil {
		utils.Logger.Fatalf("create %s fail, %v", filename, err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(YAML_TEMPLATE)); err != nil{
		utils.Logger.Fatalf("create %s fail, %v", filename, err)
	}
	return nil
}