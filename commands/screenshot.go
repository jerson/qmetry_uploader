package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"qmetry_uploader/modules/utils"
	"strings"
)

// ScreenshotOptions ...
type ScreenshotOptions struct {
	Model       string
	Case        string
	Description string
}

// ScreenshotAndroidOptions ...
type ScreenshotAndroidOptions struct {
	ScreenshotOptions
	ADB string
}

// ScreenshotAndroid
//
// initial script for bash
// MODEL="$(echo $1 | awk '{$1=$1};1')"
// CASE="$(echo $2 | awk '{$1=$1};1')"
// DESCRIPTION="$(echo $3 | awk '{$1=$1};1')"
// SLUG="$(echo ${MODEL}_${CASE}_${DESCRIPTION} | iconv -t ascii//TRANSLIT | sed -E s/[^a-zA-Z0-9_]+/-/g | sed -E s/^-+\|-+$//g)"
// FILE=$SLUG.png
// echo "capture: $FILE";
// adb exec-out screencap -p > $FILE
func ScreenshotAndroid(options ScreenshotAndroidOptions) error {

	caseName := strings.Trim(options.Case, " ")
	description, err := utils.Slug(strings.Trim(options.Description, " "))
	if err != nil {
		description = strings.Trim(options.Description, " ")

	}

	model := strings.Trim(options.Model, " ")
	name := strings.Join([]string{model, caseName, description}, "_")
	output := fmt.Sprintf("%s.png", name)

	cmd := exec.Command(options.ADB, "exec-out", "screencap", "-p")
	outfile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	stdError, _ := ioutil.ReadAll(cmdErr)
	errorString := string(stdError)
	if errorString != "" {
		log.Errorf(errorString)
	} else {
		log.Infof("output: %s\n", output)
	}

	return nil
}

// ScreenshotSample ...
func ScreenshotSample(options ScreenshotAndroidOptions) error {

	caseName := strings.Trim(options.Case, " ")
	description, err := utils.Slug(strings.Trim(options.Description, " "))
	if err != nil {
		description = strings.Trim(options.Description, " ")

	}

	model := strings.Trim(options.Model, " ")
	name := strings.Join([]string{model, caseName, description}, "_")
	output := fmt.Sprintf("%s.png", name)

	cmd := exec.Command(options.ADB, "exec-out", "screencap", "-p", ">", output)
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	stdOutput, _ := ioutil.ReadAll(cmdOut)
	stdError, _ := ioutil.ReadAll(cmdErr)

	errorString := string(stdError)
	if errorString != "" {
		log.Errorf(errorString)
	}
	log.Debug(string(stdOutput))

	log.Infof("output: %s\n", output)

	return nil
}
