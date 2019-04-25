package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"

	"github.com/urfave/cli"
)

func setup() {
	_ = config.ReadDefault()
}
func main() {

	app := cli.NewApp()
	app.Name = "Qmetry uploader"
	app.Usage = "sube facilmente tus evidencias a Qmetry"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name:    "compress",
			Aliases: []string{"c"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output dir",
				},
			},
			Usage: "compress images",
			Action: func(c *cli.Context) error {
				readContext(c)

				err := commands.Compress()

				return err
			},
		},
		{
			Name:    "report",
			Aliases: []string{"r"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output dir",
				},
			},
			Usage: "show report",
			Action: func(c *cli.Context) error {
				readContext(c)

				data, err := commands.Report()
				if data != nil {
					printJSON(data)
				}
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func readContext(c *cli.Context) {

	input := c.String("input")
	if input != "" {
		config.Vars.Dir.Input = input
	} else {
		config.Vars.Dir.Input = "./"
	}

	output := c.String("output")
	if output != "" {
		config.Vars.Dir.Output = output
	} else {
		config.Vars.Dir.Output = "./output"
	}

}

func printJSON(data interface{}) {

	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}
