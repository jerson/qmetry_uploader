package osx

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"os"
	"os/exec"
	"os/user"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/utils"
)

// OpenApp ...
func OpenApp(name string) error {

	err := exec.Command("open", "-a", name).Run()
	if err != nil {
		return err
	}
	return nil
}

// GetAutomatorDir ...
func GetAutomatorDir() (string, error) {
	output := ""
	usr, err := user.Current()
	if err != nil {
		return output, err
	}
	output = fmt.Sprintf("%s/%s", usr.HomeDir, config.Vars.Dir.Automator)
	return output, nil
}

// GetAutomatorFile ...
func GetAutomatorFile(name string) (string, error) {
	output := ""
	dir, err := GetAutomatorDir()
	if err != nil {
		return output, err
	}
	output = fmt.Sprintf("%s/%s", dir, name)
	return output, nil
}

// LoadAssets ...
func LoadAssets() error {
	box := packr.New("assets", "./assets")
	bytes, err := box.Find("automator.zip")
	if err != nil {
		return err
	}

	dir, err := GetAutomatorDir()
	if err != nil {
		return err
	}
	err = os.Mkdir(dir, 0777)
	if err != nil {
		return err
	}
	_, err = utils.Unzip(bytes, dir)
	return nil
}
