package osx

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"qmetry_uploader/modules/assets"
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
	box := assets.Load()
	bytes, err := box.Find("assets/automator.zip")
	if err != nil {
		return err
	}

	dir, err := GetAutomatorDir()
	if err != nil {
		return err
	}
	_ = os.Mkdir(dir, 0777)
	_, err = utils.Unzip(bytes, dir)
	return nil
}
