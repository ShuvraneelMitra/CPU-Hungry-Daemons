package gui

import (
	// "image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/canvas"
	// "fyne.io/fyne/v2/layout"
)

const (
	WINDOW_WIDTH = 1200
	WINDOW_HEIGHT = 615
)

func Run() {
	newApp := app.New()
	win := newApp.NewWindow("Hungry-Daemons")
	win.SetMaster()

	win.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))
	screen := newApp.Driver().AllWindows()[0]
	screen.CenterOnScreen()

	content := getLayout()

	win.SetContent(content)

	win.ShowAndRun()
}