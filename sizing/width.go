package sizing

import (
	"fmt"
	"image"

	"github.com/ddmytro-m/asciitor/internal/palette"
)

type OutputWidth interface {
	GetWidth(img image.Image, palette *palette.Palette, blockSize int) (int, error)
}

type WidthAuto struct{}

func (o WidthAuto) GetWidth(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	imageWidth := img.Bounds().Dx()
	if imageWidth <= 0 {
		return -1, fmt.Errorf("input image's width is invalid: %dpx", imageWidth)
	}
	return imageWidth, nil
}

type WidthPixels struct {
	Pixels int
}

func (o WidthPixels) GetWidth(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	if o.Pixels <= 0 {
		return -1, fmt.Errorf("invalid width size: %dpx", o.Pixels)
	}
	return o.Pixels, nil
}

type WidthCharacters struct {
	Character rune
	Amount    int
}

func (o WidthCharacters) GetWidth(img image.Image, palette *palette.Palette, blockSize int) (int, error) {
	if o.Amount <= 0 {
		return -1, fmt.Errorf("characters amount is invalid: %d", o.Amount)
	}

	width := palette.GetCharacterWidth(o.Character)
	if width <= 0 {
		return -1, fmt.Errorf("\"%c\" character's width is unknown", o.Character)
	}

	return width * o.Amount, nil
}
