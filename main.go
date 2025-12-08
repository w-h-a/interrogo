package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/interrogo/cmd"
)

func main() {
	app := &cli.App{
		Name:  "interrogo",
		Usage: "Judge your agents in CI before deploy",
		Commands: []*cli.Command{
			{
				Name: "judge",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Path to evaluator config",
						Required: true,
					},
				},
				Action: cmd.Judge,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
