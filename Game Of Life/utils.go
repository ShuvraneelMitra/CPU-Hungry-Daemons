package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

func getIcon() image.Image {
	iconFile, err := os.Open("logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer iconFile.Close()

	iconImage, err := png.Decode(iconFile)
	if err != nil {
		log.Fatal(err)
	}

	return iconImage
}
