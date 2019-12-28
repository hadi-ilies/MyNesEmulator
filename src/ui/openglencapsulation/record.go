package openglencapsulation

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/fogleman/nes/nes"
)

func saveGIF(path string, frames []image.Image) error {
	var palette []color.Color

	for _, c := range nes.Palette {
		palette = append(palette, c)
	}
	g := gif.GIF{}
	for i, src := range frames {
		if i%3 != 0 {
			continue
		}
		dst := image.NewPaletted(src.Bounds(), palette)
		draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
		g.Image = append(g.Image, dst)
		g.Delay = append(g.Delay, 5)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, &g)
}

func Record(frames []image.Image) {
	for i := 0; i < 1000; i++ {
		path := fmt.Sprintf("%03d.gif", i)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			saveGIF(path, frames)
			return
		}
	}
}
