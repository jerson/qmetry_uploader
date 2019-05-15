package commands

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	fontButton *ui.FontButton
	alignment  *ui.Combobox

	attrstr *ui.AttributedString
)

// UI ...
func UI() error {

	return ui.Main(setupUI)
}

func setupUI() {

	main := ui.NewWindow("Choose files to upload (*.apk,*.ipa,*.zip)", 640, 480, false)
	main.SetMargined(true)
	main.OnClosing(func(*ui.Window) bool {
		main.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		main.Destroy()
		return true
	})

	container := ui.NewVerticalBox()
	container.SetPadded(true)
	main.SetChild(container)

	header := setupHeader()
	container.Append(header, false)

	body := setupBody()
	container.Append(body, true)

	footer := setupFooter()
	container.Append(footer, false)

	main.Show()
}

func setupBody() *ui.Box {
	body := ui.NewHorizontalBox()
	body.SetPadded(true)

	mh := newModelHandler()
	model := ui.NewTableModel(mh)
	table := ui.NewTable(&ui.TableParams{
		Model: model,
	})

	body.Append(table, true)

	return body
}
func setupHeader() *ui.Box {
	header := ui.NewHorizontalBox()
	header.SetPadded(true)

	upDirButton := ui.NewButton("Up dir")
	header.Append(upDirButton, false)

	currentDir := ui.NewEntry()
	currentDir.SetReadOnly(true)
	header.Append(currentDir, true)

	changeDirButton := ui.NewButton("Change dir")
	header.Append(changeDirButton, false)

	return header
}
func setupFooter() *ui.Box {
	footer := ui.NewHorizontalBox()
	footer.SetPadded(true)

	spacing := ui.NewVerticalBox()
	footer.Append(spacing, true)

	chooseButton := ui.NewButton("Choose")
	footer.Append(chooseButton, false)
	cancelButton := ui.NewButton("Cancel")
	footer.Append(cancelButton, false)

	return footer
}

// Table definition
type modelHandler struct {
	rowCount      []int
	checkBoxValue []int
}

func newModelHandler() *modelHandler {
	mh := new(modelHandler)
	return mh
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return 100
}

func (mh *modelHandler) NumCols(m *ui.TableModel) int {
	return 1
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	tValue := make([]ui.TableValue, mh.NumCols(m))
	tValue[0] = ui.TableString("title")
	return tValue
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	return ui.TableString("test")
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	return
}
