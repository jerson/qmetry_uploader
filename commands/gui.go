package commands

import (
	log "github.com/sirupsen/logrus"

	ui "github.com/VladimirMarkelov/clui"
)

var currentPath = ""
var listBox *ui.ListBox
var frameOptions *ui.Frame

// GUI ...
func GUI() error {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, ui.AutoSize, 10, "Comprimir imagenes")
	//view.SetPack(1)
	view.SetMaximized(true)

	frame := ui.CreateFrame(view, ui.AutoSize, ui.AutoSize, ui.BorderAuto, 1)
	frame.SetPaddings(2, 2)

	btnChoose := ui.CreateButton(frame, ui.AutoSize, ui.AutoSize, "Elegir directorio", 1)
	btnChoose.OnClick(func(ev ui.Event) {
		dialog := ui.CreateFileSelectDialog("Elige directorio", "", "./", true, true)
		dialog.OnClose(func() {
			if dialog.Selected {
				currentPath = dialog.FilePath

				scenarios, err := ReportDir(currentPath)
				if err != nil {
					panic(err)
				}

				//listBox.Clear()
				for _, scenario := range scenarios {
					log.Info(scenario.Name)

					checkbox := ui.CreateCheckBox(
						frameOptions,
						ui.AutoSize,
						scenario.Name,
						1,
					)
					checkbox.SetActive(true)
					//listBox.AddChild(checkbox)
				}

			}

		})

	})

	frameOptions = ui.CreateFrame(view, ui.AutoSize, ui.AutoSize, ui.BorderAuto, 1)
	frameOptions.SetScrollable(true)
	frameOptions.SetTitle("opciones")
	//listBox = ui.CreateListBox(frameOptions, ui.AutoSize, 10, 1)
	frameOptions.SetPack(1)
	//ui.CreateLoginDialog("test", "hola")

	//listBox.AddItem("tesst")
	//listBox.AddItem("tesst")
	ui.MainLoop()

	return nil
}
