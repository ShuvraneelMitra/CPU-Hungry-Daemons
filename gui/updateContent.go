package gui

import (
	"runtime"
	"time"
	"strconv"

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

func updateStatus(layout *guiLayout) {
	go func(){
		for{
			fyne.Do(func(){
				numGoRoutines := runtime.NumGoroutine()
				layout.statusBar.body.Text = "Number of goroutines actively running = " + strconv.Itoa(numGoRoutines)
				layout.statusBar.body.Refresh()
			})
		}
	}()
}
