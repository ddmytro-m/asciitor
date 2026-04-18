package sizing

import (
	"fmt"
	"image"

	"github.com/ddmytro-m/asciitor/internal/palette"
)

type OutputHeight interface {
	GetHeight(img image.Image, palette *palette.Palette, blockSize int) (int, error)
}

type HeightAuto struct{}

func (o HeightAuto) GetHeight(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	imageHeight := img.Bounds().Dy()
	if imageHeight <= 0 {
		return -1, fmt.Errorf("input image's height is invalid: %dpx", imageHeight)
	}
	return imageHeight, nil
}

type HeightPixels struct {
	Pixels int
}

func (o WidthPixels) GetHeight(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	if o.Pixels <= 0 {
		return -1, fmt.Errorf("invalid height size: %dpx", o.Pixels)
	}
	return o.Pixels, nil
}

type HeightLines struct {
	Amount int
}

func (o HeightLines) GetHeight(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	if o.Amount <= 0 {
		return -1, fmt.Errorf("lines amount is invalid: %d", o.Amount)
	}

	height := palette.GetHeight()
	if height <= 0 {
		return -1, fmt.Errorf("line height is unknown")
	}

	return height * o.Amount, nil
}
