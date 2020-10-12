package main

import (
	"os"

	"porter/cmd/porter/actions"
	"porter/pkg/utils"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Porter",
		Usage:   "Porter is a new container engine",
		Version: actions.Version(),
		Action: func(c *cli.Context) error {
			utils.Logger.Info(`Porter is a new container engine, Use "help" see more.`)
			return nil
		},
	}

	app.Commands = append(app.Commands,
		&cli.Command{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   `Run a container, example: "porter run template.yaml"`,
			Action:  actions.Run,
		})

	app.Commands = append(app.Commands,
		&cli.Command{
			Name:    "template",
			Aliases: []string{"t", "temp"},
			Usage:   `Create a config template, example: "porter temp template.yaml"`,
			Action:  actions.Template,
		})

	app.Commands = append(app.Commands,
		&cli.Command{
			Name:   "init",
			Hidden: true,
			Action: actions.Init,
		})

	if err := app.Run(os.Args); err != nil {
		utils.Logger.Fatal("Porter have some trouble: %v", err)
	}
}
