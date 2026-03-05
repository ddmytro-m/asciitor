package converter

import (
	"fmt"

	"github.com/ddmytro-m/asciitor/internal/graphics"
	"github.com/ddmytro-m/asciitor/internal/palette"
)

type GlyphFeatureVector struct {
	character     rune
	featureVector graphics.FeatureVector
}

type Converter struct {
	Rows    int
	Columns int

	VerticalResolution   int
	HorizontalResolution int

	Palette *palette.Palette

	glyphFeatureVectors []GlyphFeatureVector
}

func NewConverter(rows, columns int, palette *palette.Palette, verticalResolution, horizontalResolution int) (*Converter, error) {
	if rows <= 0 || columns <= 0 {
		return nil, fmt.Errorf("invalid converted dimensions: (%dx%d)", rows, columns)
	}
	if verticalResolution <= 0 || horizontalResolution <= 0 {
		return nil, fmt.Errorf("invalid converter resolution: (%dx%d)", verticalResolution, horizontalResolution)
	}

	c := new(Converter)

	c.Rows = rows
	c.Columns = columns

	c.VerticalResolution = verticalResolution
	c.HorizontalResolution = horizontalResolution
	c.Palette = palette

	return c, nil
}

func (c *Converter) Load() error {
	featureVectors, err := c.getPaletteFeatureVectors()
	if err != nil {
		return err
	}

	c.glyphFeatureVectors = featureVectors

	return nil
}

func (c *Converter) GetMinImageSize() (w, h int) {
	return c.HorizontalResolution * c.Columns, c.VerticalResolution * c.Rows
}

func (c *Converter) GetMaximumResolution() (h, v int) {
	return c.Palette.Face.GetCharacterDimensions()
}

func (c *Converter) getPaletteFeatureVectors() ([]GlyphFeatureVector, error) {
	featureVectors := make([]GlyphFeatureVector, len(c.Palette.GlyphBitmaps))

	for i, glyphBitmap := range c.Palette.GlyphBitmaps {
		fv, err := graphics.NewFeatureVectorFromBitmap(glyphBitmap.Bitmap, c.HorizontalResolution, c.VerticalResolution)
		if err != nil {
			return nil, err
		}
		featureVectors[i] = GlyphFeatureVector{glyphBitmap.Character, *fv}
	}

	return featureVectors, nil
}

func (c *Converter) getClosestCharacter(chunk graphics.Bitmap) (rune, error) {
	chunkFeatureVector, err := graphics.NewFeatureVectorFromBitmap(chunk, c.HorizontalResolution, c.VerticalResolution)
	if err != nil {
		return 0, err
	}

	closestChar := c.glyphFeatureVectors[0].character
	closestDistance, err := chunkFeatureVector.GetDistance(&c.glyphFeatureVectors[0].featureVector)
	if err != nil {
		return 0, err
	}

	for k := 1; k < len(c.glyphFeatureVectors); k++ {
		gfv := c.glyphFeatureVectors[k]
		d, err := chunkFeatureVector.GetDistance(&gfv.featureVector)
		if err != nil {
			return 0, err
		}

		if d < closestDistance {
			closestChar = gfv.character
			closestDistance = d
		}
	}
	return closestChar, nil
}

func (c *Converter) Convert(bitmap graphics.Bitmap) ([][]rune, error) {
	if !c.Palette.IsRendered() {
		return nil, fmt.Errorf("can't convert bitmap: palette isn't rendered")
	}

	if len(c.glyphFeatureVectors) == 0 {
		return nil, fmt.Errorf("can't convert bitmap: no rendered feature vectors")
	}

	minW, minH := c.GetMinImageSize()
	if bitmap.Width < minW || bitmap.Height < minH {
		return nil, fmt.Errorf("can't convert bitmap: it is too small")
	}

	output := make([][]rune, c.Rows)

	chunks, err := bitmap.SplitIntoChunks(c.Columns, c.Rows)
	if err != nil {
		return nil, err
	}

	for i, row := range chunks {
		output[i] = make([]rune, c.Columns)
		for j, chunk := range row {
			closestChar, err := c.getClosestCharacter(chunk)
			if err != nil {
				return nil, err
			}

			output[i][j] = closestChar
		}
	}

	return output, nil
}
