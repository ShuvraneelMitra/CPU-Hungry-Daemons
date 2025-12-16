package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func getLayout() *fyne.Container {
	var header fyne.CanvasObject = func() fyne.CanvasObject {
		t := canvas.NewText(
			"hungry-daemons: an exploration of goroutines by insipidintegrator",
			theme.Color(theme.ColorNameForeground),
		)
		t.TextStyle = fyne.TextStyle{Monospace: true}
		t.TextSize = theme.TextSize() * 0.75
		return t
	}()

	var footer fyne.CanvasObject = func() fyne.CanvasObject {
		t := canvas.NewText(
			"\u00A9 ShuvraneelMitra",
			theme.Color(theme.ColorNameForeground),
		)
		t.TextStyle = fyne.TextStyle{Monospace: true}
		t.TextSize = theme.TextSize() * 0.75
		return t
	}()

	themedHeader := getThemedHeaderandFooter(header)
	themedFooter := getThemedHeaderandFooter(footer)

    sidebar := widget.NewLabel("Default")
	themedSidebar := getThemedSidebar(sidebar)
    top := widget.NewLabel("Default")
    bottom := widget.NewLabel("Default")

    main_win := container.NewVSplit(
        top,
        bottom,
    ) 
	main_win.Offset = 0.85

	split := container.NewHSplit(
		themedSidebar,
		main_win,
	)
	split.Offset = 0.15

    return container.NewBorder(
        themedHeader, 
        themedFooter, 
        nil,   
        nil,    
        split, 
    )
}
