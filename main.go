package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"

	"github.com/urfave/cli"
)

func setup() {
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	//_ = config.ReadDefault()
}
func main() {

	setup()

	app := cli.NewApp()
	app.Name = "Qmetry uploader"
	app.Version = "0.0.1"
	app.Usage = "sube facilmente tus evidencias a Qmetry"

	app.Flags = []cli.Flag{}

	app.CommandNotFound = func(c *cli.Context, name string) {
		_ = commands.GUI()
	}
	app.Commands = []cli.Command{
		{
			Name:    "merge-images",
			Aliases: []string{"m"},
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
			Usage: "merge images into merged file",
			Action: func(c *cli.Context) error {
				readContext(c)
				err := commands.MergeImages()
				return err
			},
		},
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
		{
			Name:    "screenshot-android",
			Aliases: []string{"sa"},
			Flags: []cli.Flag{
				cli.StringFlag{

					Name:  "adb, a",
					Usage: "ADB path",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name",
				},
			},
			Usage: `
screenshot for android using adb
ex:

qmetry-uploader screenshot-android J2 AMM-12112 01
qmetry-uploader screenshot-android J2 AMM-12112 02
qmetry-uploader screenshot-android J2 AMM-12112 "sample case"

`,
			Action: func(c *cli.Context) error {
				readContext(c)

				model := c.Args().Get(0)
				caseName := c.Args().Get(1)
				description := c.Args().Get(2)
				adb := c.String("adb")
				if adb == "" {
					adb = "adb"
				}

				options := commands.ScreenshotAndroidOptions{
					ScreenshotOptions: commands.ScreenshotOptions{
						Model:       model,
						Case:        caseName,
						Description: description,
					},
					ADB: adb,
				}

				return commands.ScreenshotAndroid(options)

			},
		},
		{
			Name:    "gui",
			Aliases: []string{"g"},
			Flags:   []cli.Flag{},
			Usage:   "show GUI",
			Action: func(c *cli.Context) error {

				return commands.GUI()
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
	log.Debug(string(output))
}
