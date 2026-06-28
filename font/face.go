package font

import (
	"fmt"

	"github.com/ddmytro-m/asciitor/internal/bridges/freetype"
	"github.com/ddmytro-m/asciitor/internal/graphics"
)

type Face struct {
	font  *Font
	index int

	styleName string

	isMonospace bool
}

func newFace(font *Font, properties freetype.FaceProperties) *Face {
	face := Face{font: font, index: properties.Index, styleName: properties.StyleName, isMonospace: properties.Monospace}
	return &face
}

func (f *Face) StyleName() string {
	return f.styleName
}

func (f *Face) FullName() string {
	return f.font.familyName + " " + f.styleName
}

func (f *Face) IsMonospace() bool {
	return f.isMonospace
}

func (f *Face) GetParams() freetype.FaceParams {
	return freetype.FaceParams{
		FontParams: f.font.GetParams(),
		FaceIndex:  f.index,
	}
}

type GlyphBitmap struct {
	Charcode rune
	Bitmap   graphics.Bitmap
}

func (f *Face) Render(characters []rune, fontSize int) ([]GlyphBitmap, error) {
	renderOutput, err := freetype.Render(f.GetParams(), fontSize, characters)
	if err != nil {
		return nil, fmt.Errorf("freetype error during render: %v", err)
	}

	glyphBitmaps := make([]GlyphBitmap, 0, len(characters))

	for _, char := range renderOutput.Characters {
		charCode := char.Charcode
		width := char.Advance
		height := renderOutput.TextHeight

		charBitmap, err := graphics.NewBitmap(width, height)
		if err != nil {
			fmt.Printf("failed to create bitmap for character %q: %v\n", char, err)
			continue
		}

		if char.BitmapBuffer != nil {
			for j := 0; j < char.BitmapHeight; j++ {
				targetRow := j + char.TopShift
				if targetRow < 0 || targetRow >= height {
					continue
				}

				availableWidth := width - char.LeftShift
				if availableWidth <= 0 {
					continue
				}

				copyWidth := min(char.BitmapWidth, availableWidth)
				if copyWidth <= 0 {
					continue
				}

				srcStart := j * char.BitmapWidth
				srcEnd := srcStart + copyWidth

				dstStart := (targetRow * width) + max(0, char.LeftShift)
				dstEnd := dstStart + copyWidth

				if dstEnd <= len(charBitmap.Buffer) {
					copy(charBitmap.Buffer[dstStart:dstEnd], char.BitmapBuffer[srcStart:srcEnd])
				}
			}
		}

		glyphBitmaps = append(glyphBitmaps, GlyphBitmap{
			Charcode: charCode,
			Bitmap:   *charBitmap,
		})
	}

	return glyphBitmaps, nil
}
