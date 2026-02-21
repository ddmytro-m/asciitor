package palette

import (
	"fmt"

	"github.com/ddmytro-m/asciitor/internal/bitmap"
	"github.com/ddmytro-m/asciitor/internal/font"
)

type Palette struct {
	face *font.FontFace // @TODO: fallback font

	characters         []rune
	renderedCharacters map[rune]*bitmap.Bitmap
	isRendered         bool
}

func NewPalette(face *font.FontFace, characters []rune) (*Palette, error) {
	palette := new(Palette)

	err := palette.SetFont(face)
	if err != nil {
		return nil, err
	}

	palette.characters = characters

	return palette, nil
}

func (p *Palette) SetFont(face *font.FontFace) error {
	if !face.IsLoaded() {
		return fmt.Errorf("font face is not loaded")
	} else if !face.IsMonospace {
		return fmt.Errorf("only monospace fonts are currently supported")
	}

	p.face = face
	p.isRendered = false

	return nil
}

func (p *Palette) RenderCharacters() error {
	// @TODO: characters rendering

	return fmt.Errorf("under construction")
}
