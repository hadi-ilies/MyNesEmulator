package openglencapsulation

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

//create a texture
func CreateTexture() uint32 {
	var texture uint32

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

//set a two-dimensional texture image
func SetTexture(image *image.RGBA) {
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(image.Rect.Size().X), int32(image.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(image.Pix))
}
