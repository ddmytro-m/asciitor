package options

import (
	"strconv"
	"strings"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
	"github.com/ddmytro-m/asciitor/sizing"
)

type WidthResolver interface {
	Matcher[string]
	Resolver[string, sizing.OutputWidth]
}

type widthOriginal struct{}

func (widthOriginal) Match(s string) bool {
	return strings.TrimSpace(s) == "original"
}

func (widthOriginal) Resolve(string) (sizing.OutputWidth, error) {
	return sizing.WidthAuto{}, nil
}

type widthTerminal struct{ term Terminal }

func (widthTerminal) Match(s string) bool {
	return strings.TrimSpace(s) == "tw"
}

func (w widthTerminal) Resolve(string) (sizing.OutputWidth, error) {
	if w.term.Cols <= 0 {
		return sizing.WidthAuto{}, nil
	}
	return sizing.WidthCharacters{Character: w.term.Ref, Amount: w.term.Cols}, nil
}

type widthPixels struct{}

func (widthPixels) Match(s string) bool {
	s = strings.TrimSpace(s)
	if !rePx.MatchString(s) {
		return false
	}
	_, ok := parseAmount(s[:len(s)-2])
	return ok
}

func (widthPixels) Resolve(s string) (sizing.OutputWidth, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s[:len(s)-2])
	if err != nil {
		return nil, err
	}
	return sizing.WidthPixels{Pixels: n}, nil
}

type widthCharacters struct{}

func (widthCharacters) Match(s string) bool {
	s = strings.TrimSpace(s)
	if !reCols.MatchString(s) {
		return false
	}
	_, ok := parseAmount(s[:len(s)-1])
	return ok
}

func (widthCharacters) Resolve(s string) (sizing.OutputWidth, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return nil, err
	}
	return sizing.WidthCharacters{Character: rune(s[len(s)-1]), Amount: n}, nil
}

func NewWidthChain(t Terminal) *resolve.Chain[string, sizing.OutputWidth] {
	chain := resolve.NewChain[string, sizing.OutputWidth]()
	for _, w := range []WidthResolver{
		widthOriginal{},
		widthTerminal{t},
		widthPixels{},
		widthCharacters{},
	} {
		chain.AddLink(resolve.NewNode(w, w))
	}
	return chain
}