package main

//go:generate rm -rf assets/automator.zip
//go:generate zip -r assets/automator.zip assets/automator
//go:generate packr2

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
	"qmetry_uploader/modules/terminal"
	"qmetry_uploader/modules/utils"
)

func setup() {
	log.SetLevel(log.DebugLevel)
	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}
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
					Value: "",
					Usage: "Output file",
				},
			},
			Category:    "evidences",
			Description: "merge images into one merged file",
			Usage:       "qmetry-uploader merge-images",
			UsageText: `
qmetry-uploader merge-images
qmetry-uploader merge-images -o ./output.png
qmetry-uploader merge-images --input=./images
qmetry-uploader merge-images -i ./images
qmetry-uploader merge-images --input=./images --output=./output.png
qmetry-uploader merge-images -i ./images -o ./output.png`,
			Action: mergeImagesAction,
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
			Action: compressAction,
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
			Action: reportAction,
		},
		{
			Name:    "screenshot-session",
			Aliases: []string{"ss"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "adb, a",
					Value: config.Vars.Binary.ADB,
					Usage: "ADB path used when platform=android",
				},
				cli.StringFlag{
					Name:  "automator, au",
					Value: config.Vars.Binary.Automator,
					Usage: "Automator used when platform=ios",
				},
				cli.StringFlag{
					Name:  "platform, p",
					Value: "android",
					Usage: "Platform: ios,android",
				},
			},
			Category:    "screenshot",
			Description: "start session for take many screenshots",
			Usage:       "qmetry-uploader screenshot-session",
			UsageText: `
qmetry-uploader screenshot-session
qmetry-uploader screenshot-session J2 AMM-12112`,
			Action: screenshotSessionAction,
		},
		{
			Name:    "screenshot",
			Aliases: []string{"s"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "adb, a",
					Value: config.Vars.Binary.ADB,
					Usage: "ADB path used when platform=android",
				},
				cli.StringFlag{
					Name:  "automator, au",
					Value: config.Vars.Binary.Automator,
					Usage: "Automator used when platform=ios",
				},
				cli.StringFlag{
					Name:  "platform, p",
					Value: "android",
					Usage: "Platform: ios,android",
				},
			},
			Category:    "screenshot",
			Description: "capture screenshot",
			Usage:       "qmetry-uploader screenshot J2 AMM-12112 01",
			UsageText: `
qmetry-uploader screenshot
qmetry-uploader screenshot J2 AMM-12112 01
qmetry-uploader screenshot J2 AMM-12112 02
qmetry-uploader screenshot J2 AMM-12112 "sample case"`,
			Action: screenshotAction,
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
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Filename with extension",
				},
			},
			Category:    "nexus",
			Description: "upload android or ios binaries to nexus",
			Usage:       "qmetry-uploader upload-nexus qa-10-10-2010.apk",
			UsageText: `
qmetry-uploader upload-nexus qa-10-10-2010.apk
qmetry-uploader upload-nexus qa-10-10-2010.ipa
qmetry-uploader upload-nexus qa-10-10-2010.zip`,
			Action: uploadNexusAction,
		},
		{
			Name:        "ui",
			Aliases:     []string{"g"},
			Flags:       []cli.Flag{},
			Category:    "gui",
			Description: "show GUI",
			Usage:       "qmetry-uploader gui",
			Action:      guiAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mergeImagesAction(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")

	input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)

	suggestion, err := utils.GetEvidenceSuggestion(input)
	if err != nil {
		log.Warn(err)
	}
	var names []string
	if suggestion.Model != "" {
		names = append(names, suggestion.Model)
	}
	if suggestion.Name != "" {
		names = append(names, suggestion.Name)
	}
	if len(names) < 1 {
		names = append(names, "merged")
	}

	output = prompt.Field("Output File", output, "", fmt.Sprintf("%s.png", strings.Join(names, "_")))

	options := commands.MergeImagesOptions{
		Input:      input,
		OutputFile: output,
	}
	return commands.MergeImages(options)
}

func compressAction(c *cli.Context) error {
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
}

func reportAction(c *cli.Context) error {

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
}

func screenshotSessionAction(c *cli.Context) error {

	model := c.Args().Get(0)
	caseName := c.Args().Get(1)
	adb := c.String("adb")
	platform := c.String("platform")
	automator := c.String("automator")

	if !(platform == "android" || platform == "ios") {
		return fmt.Errorf("not implemented for: %s", platform)
	}

	suggestion, err := utils.GetEvidenceSuggestion(".")
	if err != nil {
		log.Warn(err)
	}

	model = prompt.Field("Model", model, "GB|GM|GA|iOS", suggestion.Model)
	caseName = prompt.Field("Case", caseName, "AMM-000", suggestion.Name)
	adb = prompt.Field("adb path", adb, "", config.Vars.Binary.ADB)
	automator = prompt.Field("automator path", automator, "", config.Vars.Binary.Automator)

	if model == "" {
		return errors.New("missing: model")
	}
	if caseName == "" {
		return errors.New("missing: case")
	}

	var steps []string
	currentStep := 1

	commonOptions := commands.ScreenshotOptions{
		Model:       model,
		Case:        caseName,
		Description: "",
		OutputDir:   "",
	}
	if platform == "ios" {
		options := commands.ScreenshotIOSOptions{
			ScreenshotOptions: commonOptions,
			Automator:         automator,
		}
		log.Warn("preparing for screenshot, wait.... dont touch nothing please!!")
		err :=commands.ScreenshotIOSPrepare(options)
		if err != nil {
			return err
		}
		log.Info("Ready for screenshots")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("")
	fmt.Println(" Insert key for do things:")
	fmt.Println(" ------------------------")

	printHelp()

	err = terminal.InputWithoutBreakLine()
	if err != nil {
		log.Warn(err)
	}

	err = terminal.HideInput()
	if err != nil {
		log.Warn(err)
	}

	for {

		k, _, err := reader.ReadRune()
		if err != nil {
			return err
		}
		key := strings.ToUpper(string(k))

		switch key {
		case "M":
			fmt.Println("Merged images:")
			output := fmt.Sprintf("%s_%s.png", model, caseName)
			err = utils.MergeImages(steps, output)
			if err != nil {
				return err
			}
			fmt.Println("Output file: " + output)
			return nil
		case "H":
			printHelp()
			continue
		case "Q":
			fmt.Println("Quit by user action")
			return nil
		case "L":
			fmt.Println("List:")

			output, err := json.MarshalIndent(steps, "", " ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(output))
			continue
		case "D":
			if len(steps) < 1 {
				fmt.Println("Nothing to remove")
				continue
			}
			last := steps[len(steps)-1]
			fmt.Println(fmt.Sprintf("Removed last: %s", last))
			err := os.Remove(last)
			if err != nil {
				panic(err)
			}
			steps = steps[:len(steps)-1]
			currentStep--

			continue
		case "R":
			steps = []string{}
			currentStep = 1
			fmt.Println("Reseted data")

			continue
		case "C":

			output := fmt.Sprintf("%s_%s", model, caseName)
			err = os.MkdirAll(output, 0777)
			if err != nil {
				log.Warn(err)
			}

			commonOptions.OutputDir = output
			commonOptions.Description = fmt.Sprint(fmt.Sprintf("%02d", currentStep))

			var name string
			if platform == "android" {
				options := commands.ScreenshotAndroidOptions{
					ScreenshotOptions: commonOptions,
					ADB:               adb,
				}
				name, err = commands.ScreenshotAndroid(options)
			} else if platform == "ios" {
				options := commands.ScreenshotIOSOptions{
					ScreenshotOptions: commonOptions,
					Automator:         automator,
				}
				name, err = commands.ScreenshotIOS(options)
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			steps = append(steps, name)
			currentStep++
			continue
		case "\n":
			continue
		default:
			printHelp()

		}

	}

}

func screenshotAction(c *cli.Context) error {

	model := c.Args().Get(0)
	caseName := c.Args().Get(1)
	description := c.Args().Get(2)
	adb := c.String("adb")
	automator := c.String("automator")
	platform := c.String("platform")

	if !(platform == "android" || platform == "ios") {
		return fmt.Errorf("not implemented for: %s", platform)
	}

	if caseName == "" && description == "" && model != "" {
		description = model
		model = ""
		caseName = ""
	}

	suggestion, err := utils.GetEvidenceSuggestion(".")
	if err != nil {
		log.Warn(err)
	}

	model = prompt.Field("Model", model, "GB|GM|GA|iOS", suggestion.Model)
	caseName = prompt.Field("Case", caseName, "AMM-000", suggestion.Name)
	description = prompt.Field("Description", description, "", suggestion.Description)
	adb = prompt.Field("adb path", adb, "", config.Vars.Binary.ADB)
	automator = prompt.Field("automator path", automator, "", config.Vars.Binary.Automator)

	if model == "" {
		return errors.New("missing: model")
	}
	if caseName == "" {
		return errors.New("missing: case")
	}
	if description == "" {
		return errors.New("missing: description")
	}

	commonOptions := commands.ScreenshotOptions{
		Model:       model,
		Case:        caseName,
		Description: description,
		OutputDir:   ".",
	}
	if platform == "android" {
		options := commands.ScreenshotAndroidOptions{
			ScreenshotOptions: commonOptions,
			ADB:               adb,
		}
		_, err = commands.ScreenshotAndroid(options)
	} else if platform == "ios" {
		options := commands.ScreenshotIOSOptions{
			ScreenshotOptions: commonOptions,
			Automator:         automator,
		}
		err :=commands.ScreenshotIOSPrepare(options)
		if err != nil {
			return err
		}
		_, err = commands.ScreenshotIOS(options)
	}
	return err

}

func uploadNexusAction(c *cli.Context) error {

	username := c.String("username")
	password := c.String("password")
	server := c.String("server")
	project := c.String("project")
	name := c.String("name")
	file := c.Args().Get(0)

	file = prompt.File("Choose file to upload", file, "*.apk,*.ipa,*.zip", "")
	if file == "" {
		return errors.New("missing file")
	}

	log.Infof("Using file: %s", file)
	name = prompt.Field("Filename", name, "", filepath.Base(file))
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
			Name: name,
		},
		Username: username,
		Password: password,
		Project:  project,
		Server:   server,
	}
	url, err := commands.UploadNexus(options)
	if err != nil {
		return err
	}

	log.Info("Opening browser...")
	return utils.OpenBrowser(url)
}
func guiAction(c *cli.Context) error {

	return commands.GUI()
}
func printHelp() {
	fmt.Println(" help:")
	fmt.Println("\tC: capture screenshot")
	fmt.Println("\tM: merge screenshots and close")
	fmt.Println("\tD: delete last screenshot")
	fmt.Println("\tL: list captured screenshots")
	fmt.Println("\tR: reset all captured screenshots")
	fmt.Println("\tQ: quit")
	fmt.Println("\tH: print help")
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
