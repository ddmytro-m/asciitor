package options

import (
	"os"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
	"github.com/ddmytro-m/asciitor/internal/loaders/charsets"
)

type CharsetResolver interface {
	Matcher[string]
	Resolver[string, []rune]
}

type charsetPreset struct {
	presets map[string]string
}

func (c charsetPreset) Match(val string) bool {
	_, ok := c.presets[val]
	return ok
}

func (c charsetPreset) Resolve(val string) ([]rune, error) {
	return []rune(c.presets[val]), nil
}

type charsetFile struct{}

func (charsetFile) Match(s string) bool {
	return s != ""
}

func (charsetFile) Resolve(s string) ([]rune, error) {
	data, err := os.ReadFile(s)
	if err != nil {
		return nil, err
	}
	return []rune(string(data)), nil
}

func NewCharsetChain() *resolve.Chain[string, []rune] {
	chain := resolve.NewChain[string, []rune]()
	for _, r := range []CharsetResolver{
		charsetPreset{presets: charsets.All()},
		charsetFile{},
	} {
		chain.AddLink(resolve.NewNode(r, r))
	}
	return chain
}

var charsetChain = NewCharsetChain()
