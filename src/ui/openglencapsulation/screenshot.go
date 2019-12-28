package openglencapsulation

import (
	"image"
	"image/draw"
)

//copy screen and return it
func TakeScreenShot(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
