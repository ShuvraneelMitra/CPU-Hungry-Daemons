package gui

import (
	"time"

	"fyne.io/fyne/v2"
)

func updateTime(layout *guiLayout) {
	go func(){
		for{
			fyne.Do(func(){
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				layout.header.right.Text = currentTime
				layout.header.right.Refresh()
			})
		}
	}()
}


