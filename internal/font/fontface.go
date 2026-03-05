package font

import (
	"fmt"
	"os"

	"github.com/ddmytro-m/asciitor/internal/bridges/freetype"
	"github.com/ddmytro-m/asciitor/internal/graphics"
)

type FontFace struct {
	fontSize int

	fontBuffer []byte
	faceIndex  int // @TODO: multi-face fonts support

	familyName string
	styleName  string

	isMonospace        bool
	maxCharacterWidth  int
	maxCharacterHeight int

	isLoaded bool
}

func NewFontFace(fontSize int) *FontFace {
	face := FontFace{}
	face.ChangeFontSize(fontSize)
	return &face
}

func (f *FontFace) IsMonospace() bool {
	return f.isMonospace
}

func (f *FontFace) IsLoaded() bool {
	return f.isLoaded
}

// @COMBAK: only for monospace fonts
func (f *FontFace) GetCharacterDimensions() (maxWidth, maxHeight int) {
	return f.maxCharacterWidth, f.maxCharacterHeight
}

func (f *FontFace) getParams() freetype.FaceParams {
	return freetype.FaceParams{
		FontBuffer: f.fontBuffer,
		FontSize:   f.fontSize,
		FaceIndex:  f.faceIndex,
	}
}

func (f *FontFace) updateProperties(properties *freetype.FaceProperties) {
	f.familyName = properties.FamilyName
	f.styleName = properties.StyleName
	f.isMonospace = properties.Monospace
	f.maxCharacterWidth = properties.MaxCharacterWidth
	f.maxCharacterHeight = properties.MaxCharacterHeight
}

func (f *FontFace) LoadFontFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open font file: \"%s\"", file)
	}

	params := f.getParams()
	params.FontBuffer = data

	properties, err := freetype.GetFaceProperties(params)
	if err != nil {
		return fmt.Errorf("failed to get face properties: \"%s\"", err)
	}

	f.fontBuffer = data
	f.updateProperties(properties)
	f.isLoaded = true

	return nil
}

// @TODO: load from memory

func (f *FontFace) ChangeFontSize(fontSize int) error {
	if fontSize <= 0 {
		return fmt.Errorf("invalid font size: %d", fontSize)
	}

	if f.isLoaded {
		properties, err := freetype.GetFaceProperties(f.getParams())
		if err != nil {
			return fmt.Errorf("failed to get face properties: \"%s\"", err)
		}

		f.updateProperties(properties)
	}

	f.fontSize = fontSize

	return nil
}

type GlyphBitmap struct {
	Character rune
	Bitmap    graphics.Bitmap
}

func (f *FontFace) Render(characters []rune) ([]GlyphBitmap, error) {
	if !f.isLoaded {
		return nil, fmt.Errorf("font face is not loaded")
	}

	rawChars, err := freetype.RenderCharacters(f.getParams(), characters)
	if err != nil {
		return nil, err
	}

	glyphBitmaps := make([]GlyphBitmap, 0, len(characters))

	for i, rawChar := range rawChars {
		if rawChar == nil {
			continue
		}

		char := characters[i]
		width := rawChar.Advance
		height := f.maxCharacterHeight

		charBitmap, err := graphics.NewBitmap(width, height)
		if err != nil {
			fmt.Printf("failed to create bitmap for character %q: %v\n", char, err)
			continue
		}

		if rawChar.BitmapBuffer != nil {
			for j := 0; j < rawChar.BitmapHeight; j++ {
				srcStart := j * rawChar.BitmapWidth
				srcEnd := srcStart + rawChar.BitmapWidth

				dstStart := ((j + rawChar.TopShift) * width) + rawChar.LeftShift
				dstEnd := dstStart + min(width, rawChar.BitmapWidth)

				copy(charBitmap.Buffer[dstStart:dstEnd], rawChar.BitmapBuffer[srcStart:srcEnd])
			}
		}

		glyphBitmaps = append(glyphBitmaps, GlyphBitmap{
			Character: char,
			Bitmap:    *charBitmap,
		})
	}

	return glyphBitmaps, nil
}
