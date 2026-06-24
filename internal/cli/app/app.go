package app

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ddmytro-m/asciitor"
	"github.com/ddmytro-m/asciitor/font"
	"github.com/ddmytro-m/asciitor/internal/cli/options"
	"golang.org/x/term"
)

func Run(ctx context.Context, opts options.Values) error {
	in, err := options.OpenInput(opts.Input)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := options.OpenOutput(opts.Output)
	if err != nil {
		return err
	}
	defer out.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	// @TODO: replace this hard-coded mess with import strategies
	_, file, _, _ := runtime.Caller(0)
	wd := filepath.Dir(file)

	fontPath := filepath.Join(wd, "../../../test/data/fonts/DejaVuSansMono.ttf")
	font, err := font.NewFontFromFile(fontPath)
	if err != nil {
		return err
	}

	face, err := font.GetFace(0)
	if err != nil {
		return err
	}

	charsetPath := filepath.Join(wd, "../../../test/data/charsets/ascii.txt")
	charsetBytes, err := os.ReadFile(charsetPath)
	if err != nil {
		return err
	}
	charset := []rune(string(charsetBytes))

	const fontSize = 10
	const blockSize = 1

	o := asciitor.AsciitorOptions{
		Face:      face,
		FontSize:  fontSize,
		Charset:   charset,
		BlockSize: blockSize,
	}
	a, err := asciitor.NewAsciitor(o)
	if err != nil {
		return err
	}

	cols, rows := terminalSize()
	ref := ' '
	if len(charset) > 0 {
		ref = charset[0]
	}
	size, err := opts.OutputSize(options.Terminal{Cols: cols, Rows: rows, Ref: ref})
	if err != nil {
		return err
	}

	ro := asciitor.RenderOptions{
		OutputSize:      size,
		KeepProportions: opts.KeepProportions,
		Inverse:         opts.Inverse,
	}

	rendered, err := a.Render(img, ro)
	if err != nil {
		return err
	}

	if _, err := out.Write(flatten(rendered)); err != nil {
		return err
	}

	return nil
}

func flatten(art [][]rune) []byte {
	var buf bytes.Buffer
	for _, row := range art {
		buf.WriteString(string(row))
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func terminalSize() (cols, rows int) {
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return w, h
	}
	return 0, 0
}
