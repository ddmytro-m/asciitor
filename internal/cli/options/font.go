package options

import (
	"os"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
	"github.com/ddmytro-m/asciitor/internal/loaders/fonts"
)

type FontResolver interface {
	Matcher[string]
	Resolver[string, []byte]
}

type fontRepository interface {
	Has(string) bool
	Get(string) ([]byte, error)
}

type fontBundled struct {
	repository fontRepository
}

func (fb fontBundled) Match(val string) bool {
	return fb.repository.Has(val)
}

func (fb fontBundled) Resolve(val string) ([]byte, error) {
	return fb.repository.Get(val)
}

type fontFile struct{}

func (fontFile) Match(s string) bool {
	return s != ""
}

func (fontFile) Resolve(s string) ([]byte, error) {
	return os.ReadFile(s)
}

func NewFontChain() *resolve.Chain[string, []byte] {
	chain := resolve.NewChain[string, []byte]()
	for _, r := range []FontResolver{
		fontBundled{fonts.GetRepository()},
		fontFile{},
	} {
		chain.AddLink(resolve.NewNode(r, r))
	}
	return chain
}

var fontChain = NewFontChain()
