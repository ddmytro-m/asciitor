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

var monospaceConverter, proportionalConverter, brailleConverter *Converter
var monospacePalette, proportionalPalette, braillePalette *palette.Palette

var asciiCharset, brailleCharset []rune

func loadCharsetFromFile(path string) []rune {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("SETUP FAILURE: could not load charset file: %v", err)
		os.Exit(1)
	}
	str := string(data)
	return []rune(str)
}

func loadFaceFromFile(path string, faceIndex int) *font.Face {
	font, err := font.NewFontFromFile(path)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error during font loading: %v", err)
		os.Exit(1)
	}

	face, err := font.GetFace(faceIndex)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error during face loading: %v", err)
		os.Exit(1)
	}

	return face
}

func newPalette(charset []rune, face *font.Face, fontSize int) *palette.Palette {
	palette, err := palette.NewPalette(charset, face, fontSize)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error when creating palette: %v", err)
		os.Exit(1)
	}
	return palette
}

func newConverter(palette *palette.Palette, blockSize int) *Converter {
	converter, err := NewConverter(palette, blockSize)
	if err != nil {
		fmt.Printf("SETUP FAILURE: error when creating converter: %v", err)
		os.Exit(1)
	}
	return converter
}

func TestMain(m *testing.M) {
	fontSize := 14

	asciiPath := testutils.DataPath("charsets/ascii.txt")
	braillePath := testutils.DataPath("charsets/braille.txt")

	asciiCharset = loadCharsetFromFile(asciiPath)
	brailleCharset = loadCharsetFromFile(braillePath)

	monoPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	propPath := testutils.DataPath("fonts/Inter_18pt-Regular.ttf")
	symbPath := testutils.DataPath("fonts/NotoSansSymbols2-Regular.ttf")

	monoFace := loadFaceFromFile(monoPath, 0)
	propFace := loadFaceFromFile(propPath, 0)
	symbFace := loadFaceFromFile(symbPath, 0)

	monospacePalette = newPalette(asciiCharset, monoFace, fontSize)
	proportionalPalette = newPalette(asciiCharset, propFace, fontSize)
	braillePalette = newPalette(brailleCharset, symbFace, fontSize)

	blockSize := 3

	monospaceConverter = newConverter(monospacePalette, blockSize)
	proportionalConverter = newConverter(proportionalPalette, blockSize)
	brailleConverter = newConverter(braillePalette, blockSize)

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

func outputToString(output [][]rune) string {
	var outString strings.Builder
	for _, row := range output {
		outString.WriteString(string(row) + "\n")
	}
	return outString.String()
}

func TestConvert(t *testing.T) {
	// monospace
	siemensStarImg, err := testutils.PngToImage(testutils.DataPath("images/siemens_star.png"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	err = monospacePalette.Render()
	if err != nil {
		t.Fatalf("could not render palette: %v", err)
	}

	charW := monospacePalette.GetMonospaceWidth()
	charH := monospacePalette.GetHeight()

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
		monospaceConverter.SetBlockSize(blockSize)

		err = monospaceConverter.Load()
		if err != nil {
			t.Fatalf("failed to load converter: %v", err)
		}

		starChars, err := monospaceConverter.Convert(siemensStarImg, renderSettings)
		if err != nil {
			t.Fatalf("error during convertion: %v", err)
		}
		writeConverted(outputToString(starChars), fmt.Sprintf("siemens_star_block=%d.txt", blockSize))
	}

	// inverted images
	renderSettings.Inverse = true
	for i := range 5 {
		blockSize := i + 1
		monospaceConverter.SetBlockSize(blockSize)

		err = monospaceConverter.Load()
		if err != nil {
			t.Fatalf("failed to load converter: %v", err)
		}

		starChars, err := monospaceConverter.Convert(siemensStarImg, renderSettings)
		if err != nil {
			t.Fatalf("error during convertion: %v", err)
		}
		writeConverted(outputToString(starChars), fmt.Sprintf("siemens_star_inverse_block=%d.txt", blockSize))
	}

	// colored image
	renderSettings.Inverse = false
	monaLisaImg, err := testutils.JpegToImage(testutils.DataPath("images/mona_lisa.jpeg"))

	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}
	monospaceConverter.SetBlockSize(3)

	cols = 50
	rows = 35
	renderSettings.MaxWidth = cols * charW
	renderSettings.MaxHeight = rows * charH

	monoLisa, err := monospaceConverter.Convert(monaLisaImg, renderSettings)
	if err != nil {
		t.Errorf("error during convertion: %v", err)
	} else {
		writeConverted(outputToString(monoLisa), "mona_lisa_mono.txt")
	}

	// proportional
	proportionalSettings := RenderSettings{
		MaxWidth:        400,
		MaxHeight:       400,
		Inverse:         false,
		KeepProportions: true,
	}

	monaLisaProp, err := proportionalConverter.Convert(monaLisaImg, proportionalSettings)
	if err != nil {
		t.Errorf("error during convertion: %v", err)
	} else {
		writeConverted(outputToString(monaLisaProp), "mona_lisa_proportional.txt")
	}

	// braille
	monaLisaBraille, err := brailleConverter.Convert(monaLisaImg, proportionalSettings)
	if err != nil {
		t.Errorf("error during convertion: %v", err)
	} else {
		writeConverted(outputToString(monaLisaBraille), "mona_lisa_braille.txt")
	}
}
