package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

func is_valid(board [][]Cell, row, col int) bool {
	return 0 <= row && row < len(board) && 0 <= col && col < len(board[0])
}

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
