package options

import (
	"errors"
	"io"
	"os"
)

const stdioMarker = "-"

func OpenInput(name string) (io.ReadCloser, error) {
	if name == "" || name == stdioMarker {
		info, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeCharDevice != 0 {
			return nil, errors.New("no input file given and stdin is a terminal: pass a file path or pipe image data in")
		}
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(name)
}

func OpenOutput(name string) (io.WriteCloser, error) {
	if name == "" || name == stdioMarker {
		return nopWriteCloser{os.Stdout}, nil
	}
	return os.Create(name)
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
