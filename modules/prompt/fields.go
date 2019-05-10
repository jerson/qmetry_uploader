package prompt

import (
	"errors"

	ui "github.com/VladimirMarkelov/clui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

func chooseDir(output chan string, title, input string) {

	go func() {
		ui.InitLibrary()
		dialog := ui.CreateFileSelectDialog(title, "", input, true, true)
		dialog.OnClose(func() {
			selected := dialog.Selected
			path := dialog.FilePath
			defer ui.DeinitLibrary()
			if selected && path == "" {
				path = input
			}

			output <- path

		})
		ui.MainLoop()
	}()
}

func chooseFile(output chan string, title, input, fileMasks string) {

	go func() {
		ui.InitLibrary()
		dialog := ui.CreateFileSelectDialog(title, fileMasks, input, false, true)
		dialog.OnClose(func() {
			selected := dialog.Selected
			path := dialog.FilePath
			defer ui.DeinitLibrary()
			if selected && path == "" {
				path = input
			}

			output <- path
		})
		ui.MainLoop()
	}()
}
func File(name, value, fileMasks, defaultValue string) string {
	if value == "" {
		output := make(chan string)
		chooseFile(output, name, defaultValue, fileMasks)
		<-output
		value = <-output
	}
	return value
}
func Dir(name, value, defaultValue string) string {
	if value == "" {
		output := make(chan string)
		chooseDir(output, name, defaultValue)
		<-output
		value = <-output
	}
	return value
}

func Field(name, value, help, defaultValue string) string {
	if value == "" {
		prompt := &survey.Input{
			Message: name,
			Default: defaultValue,
			Help:    help,
		}
		err := survey.AskOne(prompt, &value, requiredField)
		if err != nil {
			log.Warn(err)
			return value
		}
	}
	return value
}

func PasswordField(name, value, help, defaultValue string) string {
	if value == "" {
		value = defaultValue
		prompt := &survey.Password{
			Message: name,
			Help:    help,
		}
		err := survey.AskOne(prompt, &value, requiredField)
		if err != nil {
			log.Warn(err)
			return value
		}

	}
	return value
}

func requiredField(ans interface{}) error {
	input := ans.(string)
	if len(input) < 1 {
		return errors.New("required field")
	}
	return nil
}
