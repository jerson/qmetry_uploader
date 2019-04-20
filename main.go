package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
)

func setup() {

	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}
}
func main() {

	setup()

	app := cli.NewApp()
	app.Name = "Qmetry uploader"
	app.Usage = "sube facilmente tus evidencias a Qmetry"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Config file",
			Value: "./config.toml",
		},
		cli.StringFlag{
			Name:  "images, i",
			Usage: "Images dir",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output dir",
		},
		cli.StringFlag{
			Name:  "result, r",
			Usage: "Result dir",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "sample",
			Aliases: []string{"s"},
			Flags: []cli.Flag{
			},
			Usage: "sample",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "compress",
			Aliases: []string{"c"},
			Flags: []cli.Flag{
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

func readContext(c *cli.Context){

	images := c.String("images")
	if images != "" {
		config.Vars.Dir.Images = images
	}
	output := c.String("output")
	if output != "" {
		config.Vars.Dir.Output = output
	}
	result := c.String("result")
	if images != "" {
		config.Vars.Dir.Result = result
	}

}

func printJSON(data interface{}) {

	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}
