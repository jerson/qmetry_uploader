package main

import (
	"encoding/json"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

func setup() {
	log.SetLevel(log.DebugLevel)

	//_ = config.ReadDefault()
}
func main() {

	setup()

	app := cli.NewApp()
	app.Name = "Qmetry uploader"
	app.Version = "0.0.2"
	app.Usage = "upload easily to Qmetry and more"

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
					Value: ".",
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "./output",
					Usage: "Output dir",
				},
			},
			Category:    "evidences",
			Description: "merge images into one merged file",
			Usage:       "qmetry-uploader merge-images",
			UsageText: `
qmetry-uploader merge-images
qmetry-uploader merge-images -o ./output
qmetry-uploader merge-images --input=./images
qmetry-uploader merge-images -i ./images
qmetry-uploader merge-images --input=./images --output=./output
qmetry-uploader merge-images -i ./images -o ./output`,
			Action: func(c *cli.Context) error {

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
					Value: ".",
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "./output",
					Usage: "Output dir",
				},
			},
			Category:    "evidences",
			Description: "compress images grouped by device and case",
			Usage:       "qmetry-uploader compress",
			UsageText: `
qmetry-uploader compress
qmetry-uploader compress -o ./output
qmetry-uploader compress --input=./images
qmetry-uploader compress -i ./images
qmetry-uploader compress --input=./images --output=./output
qmetry-uploader compress -i ./images -o ./output`,
			Action: func(c *cli.Context) error {

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
					Value: ".",
					Usage: "Input dir",
				},
			},
			Category:    "debug",
			Description: "show report for debug purposes",
			Usage:       "qmetry-uploader report",
			UsageText: `
qmetry-uploader report
qmetry-uploader report --input=./images
qmetry-uploader report -i ./images`,
			Action: func(c *cli.Context) error {

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
					Value: "adb",
					Usage: "ADB path",
				},
			},
			Category:    "screenshot",
			Description: "screenshot for android using adb",
			Usage:       "qmetry-uploader screenshot-android J2 AMM-12112 01",
			UsageText: `
qmetry-uploader screenshot-android J2 AMM-12112 01
qmetry-uploader screenshot-android J2 AMM-12112 02
qmetry-uploader screenshot-android J2 AMM-12112 "sample case"`,
			Action: func(c *cli.Context) error {

				model := c.Args().Get(0)
				caseName := c.Args().Get(1)
				description := c.Args().Get(2)
				adb := c.String("adb")

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
			Name:    "upload-nexus",
			Aliases: []string{"un"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Nexus username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Nexus password",
				},
				cli.StringFlag{
					Name:  "project, pr",
					Value: "mi-banco",
					Usage: "Nexus project",
				},
				cli.StringFlag{
					Name: "server, s",
					// project platform name
					Value: "http://mb-nexus.westus.cloudapp.azure.com/repository/%s-%s/builds/%s",
					Usage: "Nexus server template",
				},
			},
			Category:    "nexus",
			Description: "upload android or ios binaries to nexus",
			Usage:       "qmetry-uploader upload-nexus qa-10-10-2010.apk",
			UsageText: `
qmetry-uploader upload-nexus qa-10-10-2010.apk
qmetry-uploader upload-nexus qa-10-10-2010.ipa
qmetry-uploader upload-nexus qa-10-10-2010.zip`,
			Action: func(c *cli.Context) error {

				username := c.String("username")
				password := c.String("password")
				server := c.String("server")
				project := c.String("project")

				if username == "" {
					prompt := promptui.Prompt{
						Label:    "Username ",
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					username = result
				}

				if password == "" {
					prompt := promptui.Prompt{
						Label:    "Password ",
						Mask:     '*',
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					password = result
				}
				file := c.Args().Get(0)

				options := commands.UploadNexusOptions{
					UploadOptions: commands.UploadOptions{
						File: file,
					},
					Username: username,
					Password: password,
					Project:  project,
					Server:   server,
				}
				return commands.UploadNexus(options)

			},
		},
		{
			Name:        "ui",
			Aliases:     []string{"g"},
			Flags:       []cli.Flag{},
			Category:    "gui",
			Description: "show GUI",
			Usage:       "qmetry-uploader gui",
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

func requiredField(input string) error {
	if len(input) < 1 {
		return errors.New("Required field")
	}
	return nil
}
func readContext(c *cli.Context) {

	input := c.String("input")
	if input != "" {
		config.Vars.Dir.Input = input
	} else {
		config.Vars.Dir.Input = "."
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
