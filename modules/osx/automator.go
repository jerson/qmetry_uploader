package osx

import (
	"os/exec"
)

// OpenApp ...
func OpenApp(name string) error {

	err := exec.Command("open", "-a", name).Run()
	if err != nil {
		return err
	}
	return nil
}
