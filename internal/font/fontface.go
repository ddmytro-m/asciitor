package font

import (
	"fmt"
	"os"

	"github.com/ddmytro-m/asciitor/internal/bridges/freetype"
)

type FontFace struct {
	FontSize int

	fontBuffer []byte
	faceIndex  int

	FamilyName string
	StyleName  string

	IsMonospace        bool
	MaxCharacterWidth  int
	MaxCharacterHeight int

	isLoaded bool
}

func NewFontFace(fontSize int) *FontFace {
	face := FontFace{}
	face.ChangeFontSize(fontSize)
	return &face
}

func (f *FontFace) updateProperties(properties *freetype.FaceProperties) {
	f.FamilyName = properties.FamilyName
	f.StyleName = properties.StyleName
	f.IsMonospace = properties.Monospace
	f.MaxCharacterWidth = properties.MaxCharacterWidth
	f.MaxCharacterHeight = properties.MaxCharacterHeight
}

func (f *FontFace) LoadFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open font file: \"%s\"", file)
	}

	properties, err := freetype.GetFaceProperties(data, f.FontSize)
	if err != nil {
		return fmt.Errorf("failed to get face properties: \"%s\"", err)
	}

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
		properties, err := freetype.GetFaceProperties(f.fontBuffer, fontSize)
		if err != nil {
			return fmt.Errorf("failed to get face properties: \"%s\"", err)
		}

		f.updateProperties(properties)
	}

	f.FontSize = fontSize

	return nil
}

func (f *FontFace) IsLoaded() bool {
	return f.isLoaded
}
