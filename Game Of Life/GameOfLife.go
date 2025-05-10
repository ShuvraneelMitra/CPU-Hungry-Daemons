package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Game of Life",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Icon:      []pixel.Picture{pixel.PictureDataFromImage(getIcon())},
		VSync:     true,
		Resizable: true,
		Maximized: false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	im_draw := imdraw.New(nil)
	

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		im_draw.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
