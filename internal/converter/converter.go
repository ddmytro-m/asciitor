package converter

import (
	"fmt"
	"image"
	"math"

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
	if blockSize < 1 {
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

func (c *Converter) IsLoaded() bool {
	if c.palette.Version() != c.loadedPaletteVersion {
		c.isLoaded = false
	}

	return c.isLoaded
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

func (c *Converter) SetBlockSize(blockSize int) error {
	if blockSize < 1 {
		return fmt.Errorf("invalid block size: %dpx", blockSize)
	}

	if blockSize == c.blockSize {
		return nil
	}

	c.blockSize = blockSize
	c.isLoaded = false
	return nil
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

func (c *Converter) calculateDimensions(img image.Image, settings *RenderSettings) (width, height int) {
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

func (c *Converter) prepareBitmap(img image.Image, settings *RenderSettings) (*graphics.Bitmap, error) {
	width, height := c.calculateDimensions(img, settings)
	bitmap, err := graphics.NewBitmap(width, height)
	if err != nil {
		return nil, err
	}

	bitmap.FillWithImage(img)

	if settings.Inverse {
		bitmap.Invert()
	}

	return bitmap, nil
}

// func (c *Converter) prepareSAT(image image.Image, settings RenderSettings) (*graphics.SummedAreaTable, error) { }

func (c *Converter) convertMonospace(img image.Image, settings *RenderSettings) ([][]rune, error) {
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
			posX := x * charWidth
			posY := y * charHeight

			chunk, err := bitmap.GetChunk(posX, posY, posX+charWidth, posY+charHeight)
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

	for i := 1; i < len(c.glyphDensityMaps); i++ {
		gdm := c.glyphDensityMaps[i]
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

func (c *Converter) convertProportional(img image.Image, settings *RenderSettings) ([][]rune, error) {
	bitmap, err := c.prepareBitmap(img, settings)
	if err != nil {
		return nil, err
	}

	charHeight := c.palette.GetHeight()
	rows := bitmap.Height / charHeight

	output := make([][]rune, rows)

	for y := range rows {
		line := make([]rune, 0)
		posY := y * charHeight
		cx := 0
		for {
			closestChar := rune(-1)
			closestCharWidth := 0
			distance := math.Inf(1)

			for _, gdm := range c.glyphDensityMaps {
				if cx+gdm.width > bitmap.Width {
					// adding this character will result in width overflow
					continue
				}

				x1 := cx
				y1 := posY

				x2 := x1 + gdm.width
				y2 := y1 + charHeight

				// @TODO: cache different rectangles
				chunk, err := bitmap.GetChunk(x1, y1, x2, y2)
				if err != nil {
					return nil, err
				}

				chunkDensityMap, err := graphics.NewDensityMapFromBitmap(chunk, c.blockSize)
				if err != nil {
					return nil, err
				}

				charDistance, err := chunkDensityMap.GetDistance(&gdm.densityMap)
				if err != nil {
					return nil, err
				}

				// exponent is chosen to balance between mathematical pixel error (c = 0.5) and solving greedy algorithm bias (c = 1)
				// may be added as a parameter
				normalizedDistance := charDistance / math.Pow(float64(gdm.width), 0.75)

				if normalizedDistance < distance {
					distance = normalizedDistance
					closestChar = gdm.charcode
					closestCharWidth = gdm.width
				}
			}

			if closestChar != -1 {
				line = append(line, closestChar)
				cx += closestCharWidth
			} else {
				break
			}
		}

		output[y] = line
	}

	return output, nil
}

func (c *Converter) Convert(img image.Image, settings RenderSettings) ([][]rune, error) {
	if !c.IsLoaded() {
		if err := c.Load(); err != nil {
			return nil, err
		}
	}

	if c.palette.IsMonospace() {
		return c.convertMonospace(img, &settings)
	} else {
		return c.convertProportional(img, &settings)
	}
}
