package palette

import (
	"fmt"
	"unicode"

	"github.com/ddmytro-m/asciitor/internal/font"
	"github.com/ddmytro-m/asciitor/internal/utils"
)

type Palette struct {
	Face *font.FontFace // @TODO: fallback font

	charset      []rune
	GlyphBitmaps []font.GlyphBitmap

	isRendered bool
}

func NewPalette(charset []rune, face *font.FontFace) (*Palette, error) {
	palette := new(Palette)

	palette.SetCharset(charset)

	err := palette.SetFace(face)
	if err != nil {
		return nil, err
	}

	return palette, nil
}

func (p *Palette) IsRendered() bool {
	return p.isRendered
}

func (p *Palette) SetCharset(charset []rune) {
	var filteredCharset []rune
	for _, character := range charset {
		if unicode.IsGraphic(character) {
			filteredCharset = append(filteredCharset, character)
		}
	}
	filteredCharset = utils.RemoveDuplicates(filteredCharset)

	if len(filteredCharset) != len(charset) {
		fmt.Printf("ommited %d characters", len(charset)-len(filteredCharset))
	}

	p.charset = filteredCharset
	p.isRendered = false
}

func (p *Palette) SetFace(face *font.FontFace) error {
	if !face.IsLoaded() {
		return fmt.Errorf("font face is not loaded")
	} else if !face.IsMonospace() {
		return fmt.Errorf("only monospace fonts are currently supported")
	}

	p.Face = face
	p.isRendered = false

	return nil
}

func (p *Palette) Render() error {
	if p.Face == nil {
		return fmt.Errorf("font face is required")
	} else if !p.Face.IsLoaded() {
		return fmt.Errorf("font face is not loaded")
	}

	if len(p.charset) == 0 {
		return fmt.Errorf("characters list is empty, nothing to render")
	}

	// @TODO: asynchronous render
	rendered, err := p.Face.Render(p.charset)
	if err != nil {
		return fmt.Errorf("render error: %v", err)
	}

	p.GlyphBitmaps = rendered
	p.isRendered = true

	return nil
}
