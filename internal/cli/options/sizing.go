package options

import (
	"github.com/ddmytro-m/asciitor/sizing"
)

type Terminal struct {
	Cols int
	Rows int
	Ref  rune
}

func (v Values) OutputSize(t Terminal) (sizing.OutputSize, error) {
	width, err := NewWidthChain(t).Resolve(v.Width)
	if err != nil {
		return sizing.OutputSize{}, err
	}

	height, err := NewHeightChain(t).Resolve(v.Height)
	if err != nil {
		return sizing.OutputSize{}, err
	}

	return sizing.OutputSize{Width: width, Height: height}, nil
}