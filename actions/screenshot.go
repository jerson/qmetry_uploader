package actions

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
	"qmetry_uploader/modules/utils"
)

// Screenshot ...
func Screenshot(c *cli.Context) error {

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
		err := commands.ScreenshotIOSPrepare(options)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		_, err = commands.ScreenshotIOS(options)
	}
	return err

}
