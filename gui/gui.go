package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	WINDOW_WIDTH = 1200
	WINDOW_HEIGHT = 615
)

func Run() {
	newApp := app.New()
	win := newApp.NewWindow("Hungry-Daemons")

	win.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))
	screen := newApp.Driver().AllWindows()[0]
	screen.CenterOnScreen()

	hello := widget.NewLabel("Hello Fyne!")
	win.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome nya :)")
		}),
	))

	win.ShowAndRun()
}