// Implements the TileManager interface for gruid-sdl:
//
// type TileManager interface {
//   GetImage(gruid.Cell) image.Image
//   TileSize() gruid.Point
// }
//

package main

import (
	"image"
	"image/color"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/tiles"
	"golang.org/x/image/font/opentype"
)

type TileDrawer struct {
	drawer *tiles.Drawer
}

func inverseColor(c *image.Uniform) *image.Uniform {
	r, g, b, _ := c.RGBA()
	return image.NewUniform(color.RGBA{uint8(255 - r), uint8(255 - g), uint8(255 - b), 255})
}

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	fg_cell := c.Style.Fg
	bg_cell := c.Style.Bg
	fg := image.NewUniform(color.RGBA{uint8(fg_cell << 2), uint8(fg_cell << 4), uint8(fg_cell << 6), 255})
	bg := image.NewUniform(color.RGBA{uint8(bg_cell << 2), uint8(bg_cell << 4), uint8(bg_cell << 6), 255})
	return t.drawer.Draw(c.Rune, fg, bg)
}

func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}

func NewTileDrawer() (*TileDrawer, error) {
	t := &TileDrawer{}

	// IBM MDA
	font, err := opentype.Parse(ibm_mda)
	if err != nil {
		return nil, err
	}

	// Retrieve the font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72 * 2,
	})
	if err != nil {
		return nil, err
	}

	// Create new drawer for tiles using the face. Note that we could use
	// multiple faces (e.g. italic/bold/etc) -- in that case we would simply
	// define drawers for those as well and call the appropriate one in the
	// GetImage method.
	t.drawer, err = tiles.NewDrawer(face)
	if err != nil {
		return nil, err
	}
	return t, nil
}
