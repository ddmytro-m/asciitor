package options

import (
	"strconv"
	"strings"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
	"github.com/ddmytro-m/asciitor/sizing"
)

type HeightResolver interface {
	Matcher[string]
	Resolver[string, sizing.OutputHeight]
}

type heightOriginal struct{}

func (heightOriginal) Match(s string) bool {
	return strings.TrimSpace(s) == "original"
}

func (heightOriginal) Resolve(string) (sizing.OutputHeight, error) {
	return sizing.HeightAuto{}, nil
}

type heightTerminal struct{ term Terminal }

func (heightTerminal) Match(s string) bool {
	return strings.TrimSpace(s) == "th"
}

func (h heightTerminal) Resolve(string) (sizing.OutputHeight, error) {
	if h.term.Rows <= 1 {
		return sizing.HeightAuto{}, nil
	}
	return sizing.HeightLines{Amount: h.term.Rows - 1}, nil
}

type heightPixels struct{}

func (heightPixels) Match(s string) bool {
	s = strings.TrimSpace(s)
	if !rePx.MatchString(s) {
		return false
	}
	_, ok := parseAmount(s[:len(s)-2])
	return ok
}

func (heightPixels) Resolve(s string) (sizing.OutputHeight, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s[:len(s)-2])
	if err != nil {
		return nil, err
	}
	return sizing.HeightPixels{Pixels: n}, nil
}

type heightLines struct{}

func (heightLines) Match(s string) bool {
	s = strings.TrimSpace(s)
	if !reLines.MatchString(s) {
		return false
	}
	_, ok := parseAmount(s[:len(s)-1])
	return ok
}

func (heightLines) Resolve(s string) (sizing.OutputHeight, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return nil, err
	}
	return sizing.HeightLines{Amount: n}, nil
}

func NewHeightChain(t Terminal) *resolve.Chain[string, sizing.OutputHeight] {
	chain := resolve.NewChain[string, sizing.OutputHeight]()
	for _, h := range []HeightResolver{
		heightOriginal{},
		heightTerminal{t},
		heightPixels{},
		heightLines{},
	} {
		chain.AddLink(resolve.NewNode(h, h))
	}
	return chain
}