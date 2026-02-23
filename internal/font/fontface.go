package font

import (
	"fmt"
	"os"

	"github.com/ddmytro-m/asciitor/internal/bridges/freetype"
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

func (f *FontFace) LoadFromFile(file string) error {
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

func (f *FontFace) IsMonospace() bool {
	return f.isMonospace
}

func (f *FontFace) IsLoaded() bool {
	return f.isLoaded
}
