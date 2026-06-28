package palette

import (
	"fmt"
	"unicode"

	"github.com/ddmytro-m/asciitor/font"
	"github.com/ddmytro-m/asciitor/internal/utils"
)

type Palette struct {
	face     *font.Face // @TODO: fallback font
	fontSize int

	charset []rune

	// Rendered data
	glyphs          []font.GlyphBitmap
	characterWidths map[rune]int

	monospace      bool
	monospaceWidth int

	height int

	version    int
	isRendered bool
}

func NewPalette(charset []rune, face *font.Face, fontSize int) (*Palette, error) {
	palette := new(Palette)

	err := palette.SetCharset(charset)
	if err != nil {
		return nil, err
	}

	palette.SetFace(face)
	palette.SetFontSize(fontSize)

	return palette, nil
}

func (p *Palette) IsRendered() bool {
	return p.isRendered
}

func (p *Palette) Version() int {
	return p.version
}

func (p *Palette) SetCharset(charset []rune) error {
	var filteredCharset []rune
	for _, character := range charset {
		if unicode.IsGraphic(character) {
			filteredCharset = append(filteredCharset, character)
		}
	}
	filteredCharset = utils.RemoveDuplicates(filteredCharset)

	if len(filteredCharset) == 0 {
		return fmt.Errorf("given charset is empty or consists of non-graphical characters")
	}

	if len(filteredCharset) != len(charset) {
		fmt.Printf("omitted %d characters", len(charset)-len(filteredCharset))
	}

	p.charset = filteredCharset
	p.isRendered = false
	return nil
}

func (p *Palette) SetFace(face *font.Face) error {
	if face == nil {
		return fmt.Errorf("trying to assign nil as a font face")
	}

	if face == p.face {
		return nil
	}

	p.face = face
	p.isRendered = false
	return nil
}

func (p *Palette) SetFontSize(fontSize int) error {
	if fontSize < 1 {
		return fmt.Errorf("font size is too small: %dpx", fontSize)
	}

	if fontSize == p.fontSize {
		return nil
	}

	p.fontSize = fontSize
	p.isRendered = false
	return nil
}

func (p *Palette) Render() error {
	if p.face == nil {
		p.isRendered = false
		return fmt.Errorf("font face is required")
	} else if p.fontSize <= 0 {
		p.isRendered = false
		return fmt.Errorf("font size invalid")
	}

	if len(p.charset) == 0 {
		p.isRendered = false
		return fmt.Errorf("characters list is empty, nothing to render")
	}

	// @TODO: asynchronous render
	rendered, err := p.face.Render(p.charset, p.fontSize)
	if err != nil {
		p.isRendered = false
		return err
	} else if len(rendered) == 0 {
		p.isRendered = false
		return fmt.Errorf("palette render produced 0 characters (check font file and charset)")
	}

	firstHeight := rendered[0].Bitmap.Height
	p.height = firstHeight

	p.characterWidths = make(map[rune]int)

	firstCharacter := rendered[0].Charcode
	firstWidth := rendered[0].Bitmap.Width

	p.characterWidths[firstCharacter] = firstWidth

	p.monospace = true

	for _, g := range rendered {
		c := g.Charcode
		w := g.Bitmap.Width
		p.characterWidths[c] = w

		if w != firstWidth {
			p.monospace = false
		}
	}

	if p.monospace {
		p.monospaceWidth = firstWidth
	}

	p.glyphs = rendered
	p.isRendered = true
	p.version++

	return nil
}

func (p *Palette) GetGlyphBitmaps() ([]font.GlyphBitmap, error) {
	if p.isRendered {
		return p.glyphs, nil
	}

	err := p.Render()
	if err != nil {
		return nil, err
	}

	return p.glyphs, nil
}

func (p *Palette) GetCharacterWidth(character rune) int {
	return p.characterWidths[character]
}

func (p *Palette) GetMonospaceWidth() int {
	if p.monospace {
		return p.monospaceWidth
	}

	return -1
}

func (p *Palette) GetHeight() int {
	return p.height
}

func (p *Palette) IsMonospace() bool {
	return p.monospace
}
