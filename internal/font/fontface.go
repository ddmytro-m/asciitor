package font

import (
	"fmt"
	"os"
)

type FontFace struct {
	FamilyName string
	StyleName  string

	FontSize int

	IsMonospace        bool
	MaxCharacterWidth  int
	MaxCharacterHeight int

	fontBuffer []byte
	faceIndex  int

	isLoaded bool
}

func NewFontFace() *FontFace {
	return &FontFace{}
}

func (f *FontFace) LoadFromFile(file string) error {
	_, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open font file: \"%s\"", file)
	}

	// add freetype bindings

	return fmt.Errorf("under construction")
}

func (f *FontFace) IsLoaded() bool {
	return f.isLoaded
}
