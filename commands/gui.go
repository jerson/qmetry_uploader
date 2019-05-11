package commands

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

// GUI ...
func GUI() error {
	app := app.New()

	w := app.NewWindow("Qmetry uploader")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
		widget.NewScrollContainer(widget.NewVBox(
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
			widget.NewLabel("Hello Fyne!"),
		)),
	))

	w.ShowAndRun()
	return nil
}
