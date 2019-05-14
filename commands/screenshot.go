package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/modules/osx"
	"qmetry_uploader/modules/utils"
)

// ScreenshotOptions ...
type ScreenshotOptions struct {
	Model       string
	Case        string
	Description string
	OutputDir   string
}

// ScreenshotAndroidOptions ...
type ScreenshotAndroidOptions struct {
	ScreenshotOptions
	ADB string
}

// ScreenshotIOSOptions ...
type ScreenshotIOSOptions struct {
	ScreenshotOptions
	Automator string
}

// GetNameByOptions ...
func GetNameByOptions(options ScreenshotOptions) string {
	caseName := strings.Trim(options.Case, " ")
	description, err := utils.Slug(strings.Trim(options.Description, " "))
	if err != nil {
		description = strings.Trim(options.Description, " ")

	}

	model := strings.Trim(options.Model, " ")
	name := strings.Join([]string{model, caseName, description}, "_")
	return fmt.Sprintf("%s/%s.png", options.OutputDir, name)

}

// ScreenshotAndroid ...
//
// initial script for bash
// MODEL="$(echo $1 | awk '{$1=$1};1')"
// CASE="$(echo $2 | awk '{$1=$1};1')"
// DESCRIPTION="$(echo $3 | awk '{$1=$1};1')"
// SLUG="$(echo ${MODEL}_${CASE}_${DESCRIPTION} | iconv -t ascii//TRANSLIT | sed -E s/[^a-zA-Z0-9_]+/-/g | sed -E s/^-+\|-+$//g)"
// FILE=$SLUG.png
// echo "capture: $FILE";
// adb exec-out screencap -p > $FILE
func ScreenshotAndroid(options ScreenshotAndroidOptions) (string, error) {

	output := GetNameByOptions(options.ScreenshotOptions)

	cmd := exec.Command(options.ADB, "exec-out", "screencap", "-p")
	outfile, err := os.Create(output)
	if err != nil {
		return output, err
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return output, err
	}

	err = cmd.Start()
	if err != nil {
		return output, err
	}

	stdError, _ := ioutil.ReadAll(cmdErr)
	errorString := string(stdError)
	if errorString != "" {
		defer os.Remove(output)
		return output, errors.New(errorString)
	}

	log.Infof("new screenshot: %s\n", output)
	return output, nil
}

// ScreenshotIOSPrepare ...
func ScreenshotIOSPrepare(options ScreenshotIOSOptions) error {

	log.Warn("preparing for screenshot, wait.... dont touch nothing please!!")
	err := osx.OpenApp("Xcode")
	if err != nil {
		return err
	}
	defer osx.OpenApp("Terminal")
	prepareScreenShotScript, err := osx.GetAutomatorFile("prepare-screenshot.workflow")
	if err != nil {
		return err
	}
	cmd := exec.Command(options.Automator, prepareScreenShotScript)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		defer osx.OpenApp("System Preferences")
		log.Error("Open: System Preferences > Security & Privacy > Privacy > Accesibility > [enable] Terminal")
		return err
	}

	log.Info("Ready for screenshots")
	return nil
}

// ScreenshotIOS ...
func ScreenshotIOS(options ScreenshotIOSOptions) (string, error) {

	output := GetNameByOptions(options.ScreenshotOptions)

	usr, err := user.Current()
	if err != nil {
		return output, err
	}

	takeScreenShotScript, err := osx.GetAutomatorFile("take-screenshot.workflow")
	if err != nil {
		return output, err
	}

	cmd := exec.Command(options.Automator, takeScreenShotScript)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Error("Please connect device")
		return output, err
	}
	err = osx.OpenApp("Terminal")
	if err != nil {
		return output, err
	}
	log.Debug("looking for screenshot...")
	time.Sleep(4 * time.Second)

	cmd = exec.Command("bash", "-c", `ls -t ~/Desktop | grep ".png" | head -1`)
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return output, err
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return output, err
	}

	err = cmd.Start()
	if err != nil {
		return output, err
	}

	stdOutput, _ := ioutil.ReadAll(cmdOut)
	stdError, _ := ioutil.ReadAll(cmdErr)
	name := strings.TrimSpace(strings.Trim(string(stdOutput), "\n"))
	if name == "" {
		return output, errors.New("file not found")
	}

	currentScreenshot := fmt.Sprintf("%s/Desktop/%s", usr.HomeDir, name)
	errorString := string(stdError)
	if errorString != "" {
		defer os.Remove(currentScreenshot)
		return output, errors.New(errorString)
	}

	err = os.Rename(currentScreenshot, output)
	if err != nil {
		defer os.Remove(currentScreenshot)
		return output, err
	}

	log.Infof("new screenshot: %s\n", output)
	return output, nil
}
