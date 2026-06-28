package options

import (
	"errors"
	"io"
	"os"

	"net/http"
	"net/url"

	"github.com/ddmytro-m/asciitor/internal/cli/resolve"
)

type InputResolver interface {
	Matcher[string]
	Resolver[string, io.ReadCloser]
}

type stdinInput struct{}

func (stdinInput) Match(s string) bool {
	return s == "" || s == stdioMarker
}

func (stdinInput) Resolve(string) (io.ReadCloser, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return nil, errors.New("no input file given and stdin is a terminal: pass a file path or pipe image data in")
	}
	return io.NopCloser(os.Stdin), nil
}

type uriInput struct{}

func (uriInput) Match(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err == nil {
		return true
	}

	return false
}

func (uriInput) Resolve(s string) (io.ReadCloser, error) {
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type fileInput struct{}

func (fileInput) Match(s string) bool {
	return s != "" && s != stdioMarker
}

func (fileInput) Resolve(s string) (io.ReadCloser, error) {
	return os.Open(s)
}

func NewInputChain() *resolve.Chain[string, io.ReadCloser] {
	chain := resolve.NewChain[string, io.ReadCloser]()
	for _, i := range []InputResolver{stdinInput{}, uriInput{}, fileInput{}} {
		chain.AddLink(resolve.NewNode(i, i))
	}
	return chain
}

var inputChain = NewInputChain()
