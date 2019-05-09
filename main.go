package main

import (
	"encoding/json"
	"errors"
	"os"

	ui "github.com/VladimirMarkelov/clui"
	log "github.com/sirupsen/logrus"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
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
				input := c.String("input")
				output := c.String("output")
				options := commands.MergeImagesOptions{
					Input:  input,
					Output: output,
				}
				err := commands.MergeImages(options)
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
				input := c.String("input")
				output := c.String("output")
				options := commands.CompressOptions{
					Input:  input,
					Output: output,
				}
				err := commands.Compress(options)
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

				input := c.String("input")
				options := commands.ReportOptions{
					Input: input,
				}

				data, err := commands.Report(options)
				if data != nil {
					printJSON(data)
				}
				return err
			},
		},
		{
			// Agregar que el case y device sirva para crear la carpeta
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

				if model == "" {
					prompt := promptui.Prompt{
						Label:    "Model (GB|GM|GA|iOS) ",
						Default:  "",
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					model = result
				}
				if caseName == "" {
					prompt := promptui.Prompt{
						Label:    "Case (AMM-000) ",
						Default:  "",
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					caseName = result
				}
				if description == "" {
					prompt := promptui.Prompt{
						Label:    "Description ",
						Default:  "",
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					description = result
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

				if project == "" {
					prompt := promptui.Prompt{
						Label:    "Project ",
						Validate: requiredField,
					}
					result, err := prompt.Run()
					if err != nil {
						return err
					}
					project = result
				}
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

func chooseDir(input string) {
	output := make(chan string)
	dialog := ui.CreateFileSelectDialog("Choose dir", "", input, true, true)
	dialog.OnClose(func() {
		if dialog.Selected {
			output <- dialog.FilePath
		}
	})
}

func requiredField(input string) error {
	if len(input) < 1 {
		return errors.New("Required field")
	}
	return nil
}

func printJSON(data interface{}) {

	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	log.Debug(string(output))
}
