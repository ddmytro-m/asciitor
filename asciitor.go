package asciitor

import (
	"fmt"
	"image"

	"github.com/ddmytro-m/asciitor/font"
	"github.com/ddmytro-m/asciitor/internal/converter"
	"github.com/ddmytro-m/asciitor/internal/palette"
	"github.com/ddmytro-m/asciitor/sizing"
)

type Asciitor struct {
	converter *converter.Converter
	palette   *palette.Palette
	blockSize int
}

type AsciitorOptions struct {
	Face      *font.Face
	FontSize  int
	Charset   []rune
	BlockSize int
}

func NewAsciitor(options AsciitorOptions) (*Asciitor, error) {
	if options.Face == nil {
		return nil, fmt.Errorf("font face is required")
	}

	if len(options.Charset) == 0 {
		return nil, fmt.Errorf("charset is empty")
	}

	palette, err := palette.NewPalette(options.Charset, options.Face, options.FontSize)
	if err != nil {
		return nil, err
	}

	converter, err := converter.NewConverter(palette, options.BlockSize)
	if err != nil {
		return nil, err
	}

	return &Asciitor{
		converter: converter,
		palette:   palette,
		blockSize: options.BlockSize,
	}, nil
}

type RenderOptions struct {
	OutputSize sizing.OutputSize

	KeepProportions bool
	Inverse         bool
}

func (a *Asciitor) getConverterRenderSettings(img image.Image, options RenderOptions) (converter.RenderSettings, error) {
	width, err := options.OutputSize.Width.GetWidth(img, a.palette, a.blockSize)
	if err != nil {
		return converter.RenderSettings{}, fmt.Errorf("sizing error: %w", err)
	}
	height, err := options.OutputSize.Height.GetHeight(img, a.palette, a.blockSize)
	if err != nil {
		return converter.RenderSettings{}, fmt.Errorf("sizing error: %w", err)
	}

	return converter.RenderSettings{
		MaxWidth:        width,
		MaxHeight:       height,
		KeepProportions: options.KeepProportions,
		Inverse:         options.Inverse,
	}, nil
}

func (a *Asciitor) Prepare() error {
	if !a.palette.IsRendered() {
		err := a.palette.Render()
		if err != nil {
			return err
		}
	}

	if !a.converter.IsLoaded() {
		err := a.converter.Load()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Asciitor) Render(img image.Image, options RenderOptions) ([][]rune, error) {
	err := a.Prepare()
	if err != nil {
		return nil, fmt.Errorf("preparation error: %w", err)
	}

	settings, err := a.getConverterRenderSettings(img, options)
	if err != nil {
		return nil, err
	}

	output, err := a.converter.Convert(img, settings)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (a *Asciitor) SetCharset(charset []rune) error {
	if a.palette == nil {
		return fmt.Errorf("asciitor not properly initialized: palette is nil")
	}
	return a.palette.SetCharset(charset)
}

func (a *Asciitor) SetFace(face *font.Face) error {
	if a.palette == nil {
		return fmt.Errorf("asciitor not properly initialized: palette is nil")
	}
	return a.palette.SetFace(face)
}

func (a *Asciitor) SetFontSize(size int) error {
	if a.palette == nil {
		return fmt.Errorf("asciitor not properly initialized: palette is nil")
	}
	return a.palette.SetFontSize(size)
}

func (a *Asciitor) SetBlockSize(size int) error {
	if a.palette == nil {
		return fmt.Errorf("asciitor not properly initialized: palette is nil")
	}

	err := a.converter.SetBlockSize(size)
	if err != nil {
		return err
	}

	a.blockSize = size
	return nil
}
