package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
	"qmetry_uploader/modules/utils"
)

func setup() {
	log.SetLevel(log.DebugLevel)
	_ = config.ReadDefault()
}
func main() {

	setup()

	app := cli.NewApp()
	app.Name = "Qmetry uploader"
	app.Version = config.Vars.Version
	app.Usage = "Upload easily to Qmetry and more"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name:    "merge-images",
			Aliases: []string{"m"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: config.Vars.Dir.Input,
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: config.Vars.Dir.Output,
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

				input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)
				output = prompt.Dir("Output Dir", output, config.Vars.Dir.Output)

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
					Value: config.Vars.Dir.Input,
					Usage: "Input dir",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: config.Vars.Dir.Output,
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

				input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)
				output = prompt.Dir("Output Dir", output, config.Vars.Dir.Output)

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
					Value: config.Vars.Dir.Input,
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
				input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)

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
			Name:    "screenshot-session-android",
			Aliases: []string{"ssa"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "adb, a",
					Value: config.Vars.Binary.ADB,
					Usage: "ADB path",
				},
			},
			Category:    "screenshot",
			Description: "screenshot session for android using adb",
			Usage:       "qmetry-uploader screenshot-session-android",
			UsageText: `
qmetry-uploader screenshot-session-android
qmetry-uploader screenshot-session-android J2 AMM-12112`,
			Action: func(c *cli.Context) error {

				model := c.Args().Get(0)
				caseName := c.Args().Get(1)
				adb := c.String("adb")

				suggestion, err := utils.GetEvidenceSuggestion()
				if err != nil {
					log.Warn(err)
				}

				model = prompt.Field("Model", model, "GB|GM|GA|iOS", suggestion.Model)
				caseName = prompt.Field("Case", caseName, "AMM-000", suggestion.Name)
				adb = prompt.Field("adb path", adb, "", config.Vars.Binary.ADB)

				if model == "" {
					return errors.New("missing: model")
				}
				if caseName == "" {
					return errors.New("missing: case")
				}

				var steps []string
				currentStep := 1
				reader := bufio.NewReader(os.Stdin)
				fmt.Println("")
				fmt.Println(" Insert key for do things:")
				fmt.Println(" ------------------------")

				printHelp()

				for {
					k, _, err := reader.ReadRune()
					if err != nil {
						return err
					}
					key := k
					if key == 'M' {
						fmt.Println("Merged images:")
						output := fmt.Sprintf("%s_%s.png", model, caseName)
						err = utils.MergeImages(steps, output)
						if err != nil {
							return err
						}
						fmt.Println("Output file: " + output)
						return nil
					} else if key == 'P' {
						fmt.Println("Report:")

						output, err := json.MarshalIndent(steps, "", " ")
						if err != nil {
							panic(err)
						}
						fmt.Println(string(output))
						continue
					} else if key == 'D' {
						fmt.Println("Removed last")

						steps = steps[:len(steps)-1]

						continue
					} else if key == 'S' {
						options := commands.ScreenshotAndroidOptions{
							ScreenshotOptions: commands.ScreenshotOptions{
								Model:       model,
								Case:        caseName,
								Description: fmt.Sprint(currentStep),
							},
							ADB: adb,
						}

						name, err := commands.ScreenshotAndroid(options)
						if err != nil {
							return err
						}
						steps = append(steps, name)
						currentStep++
						continue
					} else if key == '\n' {

						continue

					} else {

						printHelp()

					}

				}

			},
		},
		{
			Name:    "screenshot-android",
			Aliases: []string{"sa"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "adb, a",
					Value: config.Vars.Binary.ADB,
					Usage: "ADB path",
				},
			},
			Category:    "screenshot",
			Description: "screenshot for android using adb",
			Usage:       "qmetry-uploader screenshot-android J2 AMM-12112 01",
			UsageText: `
qmetry-uploader screenshot-android
qmetry-uploader screenshot-android J2 AMM-12112 01
qmetry-uploader screenshot-android J2 AMM-12112 02
qmetry-uploader screenshot-android J2 AMM-12112 "sample case"`,
			Action: func(c *cli.Context) error {

				model := c.Args().Get(0)
				caseName := c.Args().Get(1)
				description := c.Args().Get(2)
				adb := c.String("adb")

				if caseName == "" && description == "" && model != "" {
					description = model
					model = ""
					caseName = ""
				}

				suggestion, err := utils.GetEvidenceSuggestion()
				if err != nil {
					log.Warn(err)
				}

				model = prompt.Field("Model", model, "GB|GM|GA|iOS", suggestion.Model)
				caseName = prompt.Field("Case", caseName, "AMM-000", suggestion.Name)
				description = prompt.Field("Description", description, "", suggestion.Description)
				adb = prompt.Field("adb path", adb, "", config.Vars.Binary.ADB)

				if model == "" {
					return errors.New("missing: model")
				}
				if caseName == "" {
					return errors.New("missing: case")
				}
				if description == "" {
					return errors.New("missing: description")
				}

				options := commands.ScreenshotAndroidOptions{
					ScreenshotOptions: commands.ScreenshotOptions{
						Model:       model,
						Case:        caseName,
						Description: description,
					},
					ADB: adb,
				}
				_, err = commands.ScreenshotAndroid(options)
				return err

			},
		},
		{
			Name:    "upload-nexus",
			Aliases: []string{"un"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Value: config.Vars.Nexus.Username,
					Usage: "Nexus username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: config.Vars.Nexus.Password,
					Usage: "Nexus password",
				},
				cli.StringFlag{
					Name:  "project, pr",
					Value: config.Vars.Nexus.Project,
					Usage: "Nexus project",
				},
				cli.StringFlag{
					Name:  "server, s",
					Value: config.Vars.Nexus.Server,
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
				file := c.Args().Get(0)

				file = prompt.File("Choose file to upload", file, "*.apk,*.ipa,*.zip", "")
				if file == "" {
					return errors.New("missing file")
				}

				username = prompt.Field("Username", username, "", "")
				password = prompt.PasswordField("Password", password, "", "")
				project = prompt.Field("Project", project, "", "")
				server = prompt.Field("Server", server, "", config.Vars.Nexus.Server)

				if username == "" || password == "" || project == "" || server == "" {
					return errors.New("missing fields")
				}

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

func printHelp() {
	fmt.Println(" help:")
	fmt.Println("\tS: take screenshot")
	fmt.Println("\tM: merge images and return")
	fmt.Println("\tD: delete last screenshot")
	fmt.Println("\tP: print screenshots")
	fmt.Println("")
	fmt.Println("\tNote: keys are case sensitive")
	fmt.Println("")
}
func printJSON(data interface{}) {

	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	log.Debug(string(output))
}
