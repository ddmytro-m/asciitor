package options

import (
	"io"
	"os"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
)

type OutputResolver interface {
	Matcher[string]
	Resolver[string, io.WriteCloser]
}

type stdoutOutput struct{}

func (stdoutOutput) Match(s string) bool {
	return s == "" || s == stdioMarker
}

func (stdoutOutput) Resolve(string) (io.WriteCloser, error) {
	return nopWriteCloser{os.Stdout}, nil
}

type fileOutput struct{}

func (fileOutput) Match(s string) bool {
	return s != "" && s != stdioMarker
}

func (fileOutput) Resolve(s string) (io.WriteCloser, error) {
	return os.Create(s)
}

func NewOutputChain() *resolve.Chain[string, io.WriteCloser] {
	chain := resolve.NewChain[string, io.WriteCloser]()
	for _, o := range []OutputResolver{stdoutOutput{}, fileOutput{}} {
		chain.AddLink(resolve.NewNode(o, o))
	}
	return chain
}

var outputChain = NewOutputChain()
