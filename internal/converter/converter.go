package converter

import (
	"fmt"
	"image"

	"github.com/ddmytro-m/asciitor/internal/graphics"
	"github.com/ddmytro-m/asciitor/internal/palette"
)

type GlyphDensityMap struct {
	charcode rune

	width  int
	height int

	densityMap graphics.DensityMap
}

type Converter struct {
	palette              *palette.Palette
	loadedPaletteVersion int

	blockSize        int
	glyphDensityMaps []GlyphDensityMap

	isLoaded bool
}

type RenderSettings struct {
	MaxWidth  int
	MaxHeight int

	KeepProportions bool

	Inverse bool
}

func NewConverter(pal *palette.Palette, blockSize int) (*Converter, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid block size")
	}

	if pal == nil {
		return nil, fmt.Errorf("palette cannot be nil")
	}

	return &Converter{
		palette:   pal,
		blockSize: blockSize,
	}, nil
}

func (c *Converter) Load() error {
	paletteGlyphBitmaps, err := c.palette.GetGlyphBitmaps()
	if err != nil {
		return fmt.Errorf("failed to render palette for loading: %v", err)
	}

	if len(paletteGlyphBitmaps) == 0 {
		return fmt.Errorf("cannot load converter: palette rendered 0 glyphs")
	}

	densityMaps, err := c.getDensityMaps()
	if err != nil {
		return err
	}

	c.glyphDensityMaps = densityMaps
	c.loadedPaletteVersion = c.palette.Version()
	c.isLoaded = true

	return nil
}

func (c *Converter) SetBlockSize(blockSize int) {
	if blockSize == c.blockSize {
		return
	}

	c.blockSize = blockSize
	c.isLoaded = false
}

func (c *Converter) getDensityMaps() ([]GlyphDensityMap, error) {
	paletteGlyphBitmaps, err := c.palette.GetGlyphBitmaps()
	if err != nil {
		return nil, err
	}

	densityMaps := make([]GlyphDensityMap, len(paletteGlyphBitmaps))

	for i, glyphBitmap := range paletteGlyphBitmaps {
		dm, err := graphics.NewDensityMapFromBitmap(&glyphBitmap.Bitmap, c.blockSize)
		if err != nil {
			return nil, err
		}

		densityMaps[i] = GlyphDensityMap{
			charcode:   glyphBitmap.Charcode,
			width:      glyphBitmap.Bitmap.Width,
			height:     glyphBitmap.Bitmap.Height,
			densityMap: *dm,
		}
	}

	return densityMaps, nil
}

func (c *Converter) calculateDimensions(img image.Image, settings RenderSettings) (width, height int) {
	srcW := float64(img.Bounds().Dx())
	srcH := float64(img.Bounds().Dy())
	imgRatio := srcW / srcH
	charH := c.palette.GetHeight()

	if c.palette.IsMonospace() {
		charW := c.palette.GetMonospaceWidth()
		var targetW, targetH int

		if !settings.KeepProportions {
			targetW = (settings.MaxWidth / charW) * charW
			targetH = (settings.MaxHeight / charH) * charH
		} else {
			tempW := settings.MaxWidth
			tempH := int(float64(tempW) / imgRatio)

			if tempH > settings.MaxHeight {
				tempH = settings.MaxHeight
				tempW = int(float64(tempH) * imgRatio)
			}

			targetW = (tempW / charW) * charW
			targetH = (tempH / charH) * charH
		}

		if targetW < charW {
			targetW = charW
		}
		if targetH < charH {
			targetH = charH
		}

		width, height = targetW, targetH

	} else {
		if !settings.KeepProportions {
			width = settings.MaxWidth
			height = (settings.MaxHeight / charH) * charH
		} else {
			tempW := settings.MaxWidth
			tempH := int(float64(tempW) / imgRatio)

			if tempH > settings.MaxHeight {
				tempH = settings.MaxHeight
				tempW = int(float64(tempH) * imgRatio)
			}

			height = (tempH / charH) * charH
			width = int(float64(height) * imgRatio)
		}

		if height < charH {
			height = charH
		}
		if width < 1 {
			width = 1
		}
	}

	return width, height
}

func (c *Converter) prepareBitmap(img image.Image, settings RenderSettings) (*graphics.Bitmap, error) {
	width, height := c.calculateDimensions(img, settings)
	bitmap, err := graphics.NewBitmap(width, height)
	if err != nil {
		return nil, err
	}

	bitmap.FillWithImage(img)

	return bitmap, nil
}

// func (c *Converter) prepareSAT(image image.Image, settings RenderSettings) (*graphics.SummedAreaTable, error) {

// }

func (c *Converter) convertMonospace(img image.Image, settings RenderSettings) ([][]rune, error) {
	if !c.isLoaded {
		return nil, fmt.Errorf("can't convert bitmap: palette isn't rendered")
	}

	bitmap, err := c.prepareBitmap(img, settings)
	if err != nil {
		return nil, err
	}
	// @TODO: replace bitmap with a summedareatable

	charWidth := c.palette.GetMonospaceWidth()
	charHeight := c.palette.GetHeight()

	cols := bitmap.Width / charWidth
	rows := bitmap.Height / charHeight

	output := make([][]rune, rows)

	for y := range rows {
		output[y] = make([]rune, cols)
		for x := range cols {
			imgCol := x * charWidth
			imgRow := y * charHeight

			chunk, err := bitmap.GetChunk(imgCol, imgRow, imgCol+charWidth, imgRow+charHeight)
			if err != nil {
				return nil, err
			}

			dm, err := graphics.NewDensityMapFromBitmap(chunk, c.blockSize)
			if err != nil {
				return nil, err
			}

			char, err := c.getClosestMonospaceCharacter(dm)
			if err != nil {
				return nil, err
			}

			output[y][x] = char
		}
	}

	return output, nil
}

func (c *Converter) getClosestMonospaceCharacter(dm *graphics.DensityMap) (rune, error) {
	if len(c.glyphDensityMaps) == 0 {
		return -1, fmt.Errorf("no density maps available")
	}

	closestChar := c.glyphDensityMaps[0].charcode
	closestDistance, err := dm.GetDistance(&c.glyphDensityMaps[0].densityMap)
	if err != nil {
		return 0, err
	}

	for k := 1; k < len(c.glyphDensityMaps); k++ {
		gdm := c.glyphDensityMaps[k]
		d, err := dm.GetDistance(&gdm.densityMap)
		if err != nil {
			return 0, err
		}

		if d < closestDistance {
			closestChar = gdm.charcode
			closestDistance = d
		}
	}

	return closestChar, nil
}

func (c *Converter) convertProportional(img image.Image, settings RenderSettings) ([][]rune, error) {
	return [][]rune{}, nil
}

func (c *Converter) Convert(img image.Image, settings RenderSettings) ([][]rune, error) {
	if !c.isLoaded || c.loadedPaletteVersion != c.palette.Version() {
		if err := c.Load(); err != nil {
			return nil, err
		}
	}

	if c.palette.IsMonospace() {
		return c.convertMonospace(img, settings)
	} else {
		return c.convertProportional(img, settings)
	}
}
