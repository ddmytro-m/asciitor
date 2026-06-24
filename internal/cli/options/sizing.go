package options

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ddmytro-m/asciitor/sizing"
)

type Terminal struct {
	Cols int
	Rows int
	Ref  rune
}

func (v Values) OutputSize(t Terminal) (sizing.OutputSize, error) {
	width, err := outputWidth(strings.TrimSpace(v.Width), t)
	if err != nil {
		return sizing.OutputSize{}, err
	}

	height, err := outputHeight(strings.TrimSpace(v.Height), t)
	if err != nil {
		return sizing.OutputSize{}, err
	}

	return sizing.OutputSize{Width: width, Height: height}, nil
}

func outputWidth(s string, t Terminal) (sizing.OutputWidth, error) {
	switch {
	case s == "original":
		return sizing.WidthAuto{}, nil
	case s == "tw":
		if t.Cols <= 0 {
			return sizing.WidthAuto{}, nil
		}
		return sizing.WidthCharacters{Character: t.Ref, Amount: t.Cols}, nil
	case rePx.MatchString(s):
		n, err := strconv.Atoi(s[:len(s)-2])
		if err != nil {
			return nil, err
		}
		return sizing.WidthPixels{Pixels: n}, nil
	case reCols.MatchString(s):
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return nil, err
		}
		return sizing.WidthCharacters{Character: rune(s[len(s)-1]), Amount: n}, nil
	default:
		return nil, fmt.Errorf("invalid width %q", s)
	}
}

func outputHeight(s string, t Terminal) (sizing.OutputHeight, error) {
	switch {
	case s == "original":
		return sizing.HeightAuto{}, nil
	case s == "th":
		if t.Rows <= 1 {
			return sizing.HeightAuto{}, nil
		}
		return sizing.HeightLines{Amount: t.Rows - 1}, nil
	case rePx.MatchString(s):
		n, err := strconv.Atoi(s[:len(s)-2])
		if err != nil {
			return nil, err
		}
		return sizing.HeightPixels{Pixels: n}, nil
	case reLines.MatchString(s):
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return nil, err
		}
		return sizing.HeightLines{Amount: n}, nil
	default:
		return nil, fmt.Errorf("invalid height %q", s)
	}
}
