package actions

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
	"qmetry_uploader/modules/terminal"
	"qmetry_uploader/modules/utils"
)

// ScreenShotSession ...
func ScreenShotSession(c *cli.Context) error {

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

	commonOptions := commands.ScreenShotOptions{
		Model:       model,
		Case:        caseName,
		Description: "",
		OutputDir:   "",
	}
	if platform == "ios" || platform == "ios-simulator" {
		options := commands.ScreenShotIOSOptions{
			ScreenShotOptions: commonOptions,
			Automator:         automator,
			Simulator:         platform == "ios-simulator",
		}
		err := commands.ScreenShotIOSPrepare(options)
		if err != nil {
			return err
		}

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
		case "A":
			file := prompt.File("Choose screenshot file to add", "", "*.png", "", false)
			if file == "" {
				fmt.Println("Empty file")
				continue
			}
			steps = append(steps, file)
			fmt.Println("Added file: " + file)

			continue
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
			err := os.Remove(last)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Removed last: %s", last))
			steps = steps[:len(steps)-1]
			currentStep--

			continue
		case "R":

			for _, step := range steps {
				err := os.Remove(step)
				if err != nil {
					log.Warn(err)
				}
				fmt.Println(fmt.Sprintf("Removed: %s", step))
			}
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
				options := commands.ScreenShotAndroidOptions{
					ScreenShotOptions: commonOptions,
					ADB:               adb,
				}
				name, err = commands.ScreenShotAndroid(options)
			} else if platform == "ios" || platform == "ios-simulator" {
				options := commands.ScreenShotIOSOptions{
					ScreenShotOptions: commonOptions,
					Automator:         automator,
					Simulator:         platform == "ios-simulator",
				}
				name, err = commands.ScreenShotIOS(options)
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			steps = append(steps, name)
			currentStep++
			continue
		case "\n":
		case "\r":
			continue
		default:
			printHelp()

		}

	}

}
func printHelp() {
	fmt.Println(" help:")
	fmt.Println("\tC: capture screenshot")
	fmt.Println("\tM: merge screenshots and close")
	fmt.Println("\tD: delete last screenshot")
	fmt.Println("\tL: list captured screenshots")
	fmt.Println("\tR: reset all captured screenshots")
	fmt.Println("\tA: add custom screenshot from filesystem (beta)")
	fmt.Println("\tQ: quit")
	fmt.Println("\tH: print help")
	fmt.Println("")
	fmt.Println("\tNote: keys are case sensitive")
	fmt.Println("")
}
