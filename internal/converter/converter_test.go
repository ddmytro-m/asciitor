package converter

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ddmytro-m/asciitor/font"
	"github.com/ddmytro-m/asciitor/internal/palette"
	"github.com/ddmytro-m/asciitor/internal/testutils"
)

var c *Converter
var p *palette.Palette

var asciiCharset []rune
var brailleCharset []rune

func TestMain(m *testing.M) {
	fontSize := 14

	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	font, err := font.NewFontFromFile(fontPath)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error during font loading: %v", err)
		os.Exit(1)
	}

	face, err := font.GetFace(0)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error during face loading: %v", err)
		os.Exit(1)
	}

	asciiPath := testutils.DataPath("charsets/ascii.txt")
	asciiData, err := os.ReadFile(asciiPath)
	if err != nil {
		fmt.Printf("SETUP FAILURE: could not load charset file: %v", err)
		os.Exit(1)
	}
	asciiStr := string(asciiData)
	asciiCharset = []rune(asciiStr)

	p, err = palette.NewPalette(asciiCharset, face, fontSize)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error when creating palette: %v", err)
		os.Exit(1)
	}

	err = p.Render()
	if err != nil {
		fmt.Printf("SETUP FAILURE: could not render palette: %v", err)
		os.Exit(1)
	}

	blockSize := 1

	c, err = NewConverter(p, blockSize)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error when creating converter: %v", err)
		os.Exit(1)
	}
	err = c.Load()
	if err != nil {
		fmt.Printf("SETUP FAILURE: could not load converter: %v", err)
		os.Exit(1)
	}

	m.Run()
}

func writeConverted(data string, filename string) {
	directoryPath := testutils.TmpPath("converted")
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		fmt.Printf("unable to create directory: %s", directoryPath)
		return
	}

	filePath := path.Join(directoryPath, filename)
	err = os.WriteFile(filePath, []byte(data), os.ModePerm)
	if err != nil {
		fmt.Printf("unable to create file: %s", filePath)
		return
	}
}

func TestMonospaceConvert(t *testing.T) {
	siemensStarImg, err := testutils.PngToImage(testutils.DataPath("images/siemens_star.png"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	charW := p.GetMonospaceWidth()
	charH := p.GetHeight()

	cols := 80
	rows := 40

	renderSettings := RenderSettings{
		MaxWidth:        cols * charW,
		MaxHeight:       rows * charH,
		KeepProportions: true,
		Inverse:         false,
	}

	for i := range 5 {
		blockSize := i + 1
		c.SetBlockSize(blockSize)

		err = c.Load()
		if err != nil {
			t.Fatalf("failed to load converter: %v", err)
		}

		starChars, err := c.Convert(siemensStarImg, renderSettings)
		if err != nil {
			t.Fatalf("error during convertion: %v", err)
		}

		var starString strings.Builder
		for _, row := range starChars {
			starString.WriteString(string(row) + "\n")
		}
		writeConverted(starString.String(), fmt.Sprintf("siemens_star_block=%d.txt", blockSize))
	}

	// colored image
	monaLisaImg, err := testutils.JpegToImage(testutils.DataPath("images/mona_lisa.jpeg"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	cols = 50
	rows = 35
	renderSettings.MaxWidth = cols * charW
	renderSettings.MaxHeight = rows * charH

	monaLisaChars, err := c.Convert(monaLisaImg, renderSettings)
	if err != nil {
		t.Fatalf("error during convertion: %v", err)
	}

	var monaLisaString strings.Builder
	for _, row := range monaLisaChars {
		monaLisaString.WriteString(string(row) + "\n")
	}
	writeConverted(monaLisaString.String(), "mona_lisa.txt")
}
