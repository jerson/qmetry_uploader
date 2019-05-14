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

// ScreenShotOptions ...
type ScreenShotOptions struct {
	Model       string
	Case        string
	Description string
	OutputDir   string
}

// ScreenShotAndroidOptions ...
type ScreenShotAndroidOptions struct {
	ScreenShotOptions
	ADB string
}

// ScreenShotIOSOptions ...
type ScreenShotIOSOptions struct {
	ScreenShotOptions
	Automator string
}

// GetNameByOptions ...
func GetNameByOptions(options ScreenShotOptions) string {
	caseName := strings.Trim(options.Case, " ")
	description, err := utils.Slug(strings.Trim(options.Description, " "))
	if err != nil {
		description = strings.Trim(options.Description, " ")

	}

	model := strings.Trim(options.Model, " ")
	name := strings.Join([]string{model, caseName, description}, "_")
	return fmt.Sprintf("%s/%s.png", options.OutputDir, name)

}

// ScreenShotAndroid ...
//
// initial script for bash
// MODEL="$(echo $1 | awk '{$1=$1};1')"
// CASE="$(echo $2 | awk '{$1=$1};1')"
// DESCRIPTION="$(echo $3 | awk '{$1=$1};1')"
// SLUG="$(echo ${MODEL}_${CASE}_${DESCRIPTION} | iconv -t ascii//TRANSLIT | sed -E s/[^a-zA-Z0-9_]+/-/g | sed -E s/^-+\|-+$//g)"
// FILE=$SLUG.png
// echo "capture: $FILE";
// adb exec-out screencap -p > $FILE
func ScreenShotAndroid(options ScreenShotAndroidOptions) (string, error) {

	output := GetNameByOptions(options.ScreenShotOptions)

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

	log.Infof("new screenshot: %s", output)
	return output, nil
}

// ScreenShotIOSPrepare ...
func ScreenShotIOSPrepare(options ScreenShotIOSOptions) error {

	log.Warn("preparing for screenshot, wait.... dont touch nothing please!!")
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

// GetLastFileFrom ...
func GetLastFileFrom(dir, grep string) (string, error) {
	output := ""

	cmd := exec.Command("bash", "-c", fmt.Sprintf(`ls -t %s | grep "%s" | head -1`, dir, grep))
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

	output = fmt.Sprintf("%s/%s", dir, name)
	errorString := string(stdError)
	if errorString != "" {
		return output, errors.New(errorString)
	}

	return output, nil
}

// ScreenShotIOS ...
func ScreenShotIOS(options ScreenShotIOSOptions) (string, error) {

	output := GetNameByOptions(options.ScreenShotOptions)

	usr, err := user.Current()
	if err != nil {
		return output, err
	}
	desktopDir := fmt.Sprintf("%s/%s", usr.HomeDir, "Desktop")
	screenShotExtension := ".png"

	currentLastFile, _ := GetLastFileFrom(desktopDir, screenShotExtension)

	takeScreenShotScript, err := osx.GetAutomatorFile("take-screenshot.workflow")
	if err != nil {
		return output, err
	}

	cmd := exec.Command(options.Automator, takeScreenShotScript)
	err = cmd.Run()
	if err != nil {
		log.Error("Please connect your device")
		return output, err
	}
	err = osx.OpenApp("Terminal")
	if err != nil {
		return output, err
	}
	log.Debug("looking for screenshot...")
	i := 0
	currentScreenShot := ""
	for {
		lastFile, _ := GetLastFileFrom(desktopDir, screenShotExtension)
		if lastFile != currentLastFile {
			currentScreenShot = lastFile
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
		i++
		if i > 10 {
			return output, errors.New("screenshot not found")
		}
	}

	err = os.Rename(currentScreenShot, output)
	if err != nil {
		defer os.Remove(currentScreenShot)
		return output, err
	}

	log.Infof("new screenshot: %s", output)
	return output, nil
}
