package main

//go:generate rm -rf assets/automator.zip
//go:generate zip -r assets/automator.zip assets/automator
//go:generate packr2

import (
	"os"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"qmetry_uploader/actions"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/osx"
)

func setup() {
	log.SetLevel(log.DebugLevel)

	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}

	box := packr.New("Assets", "./assets")
	err = osx.LoadAssets(box)
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
			Action: actions.MergeImages,
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
			Action: actions.Compress,
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
			Action: actions.Report,
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
			Action: actions.ScreenshotSession,
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
			Action: actions.Screenshot,
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
			Action: actions.UploadNexus,
		},
		{
			Name:        "ui",
			Aliases:     []string{"g"},
			Flags:       []cli.Flag{},
			Category:    "gui",
			Description: "show GUI",
			Usage:       "qmetry-uploader gui",
			Action:      actions.GUI,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
