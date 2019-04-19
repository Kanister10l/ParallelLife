package server

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
)

type gifBoard struct {
	Data []byte
	X    int
	Y    int
}

//CreateGif listens for game input, then at the end creates output gif file
func CreateGif(file string, width, height, scale, delay int, input chan gifBoard) {
	f, err := os.Create(fmt.Sprintf("./%s", file))
	if err != nil {
		log.Println("Error creating gif file:", err)
		return
	}

	palette := []color.Color{color.White, color.Black}
	rect := image.Rect(0, 0, width*scale, height*scale)
	images := []*image.Paletted{}

	for data := range input {
		img := image.NewPaletted(rect, palette)
		for i := 1; i < data.Y+1; i++ {
			for j := 1; j < data.X+1; j++ {
				if data.Data[i*(data.X+2)+j] == 1 {
					for k1 := 0; k1 < scale; k1++ {
						for k2 := 0; k2 < scale; k2++ {
							img.SetColorIndex(j*scale+k2, i*scale+k1, 1)
						}
					}
				}
			}
		}
		images = append(images, img)
	}

	delays := []int{}

	for i := 0; i < len(images); i++ {
		delays = append(delays, delay)
	}

	anim := gif.GIF{Delay: delays, Image: images}

	gif.EncodeAll(f, &anim)
}
