package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/pkg/errors"
)

// Texture manages a 2D texture for OpenGL
// Give it a texture png and it will create a texture from it
type Texture struct {
	textureID   uint32
	textureAddr int32
}

// NewTexture creates a 2D texture from a file
func NewTexture(program uint32, name, file string) (t Texture, err error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return t, errors.Wrap(err, "unable to open texture file")
	}

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return t, errors.Wrap(err, "unable to decode image file")
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,             // What type of texture this is
		0,                         // what level of the mipmap you are creating (default is base 0)
		gl.RGBA,                   // format to store the texture as
		int32(rgba.Rect.Size().X), // width of the texture
		int32(rgba.Rect.Size().Y), // height of the texture
		0,                // should always be 0 (legacy and no longer used)
		gl.RGBA,          // format of the source image
		gl.UNSIGNED_BYTE, // size of each element of the input
		gl.Ptr(rgba.Pix), // pointer to the actual image
	)

	textureAddr := gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.Uniform1i(textureAddr, 0)

	t = Texture{
		textureID:   textureID,
		textureAddr: textureAddr,
	}

	return t, nil
}
