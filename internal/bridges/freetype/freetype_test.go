package freetype

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ddmytro-m/asciitor/internal/testutils"
)

var validFont, dummyFont FontParams
var validFace, invalidFace FaceParams

func TestMain(m *testing.M) {
	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		fmt.Println("SETUP FAILURE: could not load font file:", err)
		os.Exit(1)
	}
	validFont.Buffer = fontData

	dummyPath := testutils.DataPath("fonts/dummy.ttf")
	dummyData, err := os.ReadFile(dummyPath)
	if err != nil {
		fmt.Println("SETUP FAILURE: could not load dummy font file:", err)
		os.Exit(1)
	}
	dummyFont.Buffer = dummyData

	validFace.FontParams = validFont
	invalidFace.FontParams = validFont

	validFace.FaceIndex = 0
	invalidFace.FaceIndex = 1234

	m.Run()
}

func TestGetFontProperties(t *testing.T) {
	properties, err := GetFontProperties(validFont)

	if properties.FamilyName != "DejaVu Sans Mono" {
		t.Errorf("expected FamilyName to be \"DejaVu Sans Mono\", got: %s", properties.FamilyName)
	}
	if properties.FacesAmount <= 0 {
		t.Errorf("expected FacesAmount to be greater than 0, got: %d", properties.FacesAmount)
	}

	// dummy test
	_, err = GetFontProperties(dummyFont)
	if err == nil {
		t.Errorf("expected error with invalid font file")
	}
}

func TestGetFaceProperties(t *testing.T) {
	// positive test
	properties, err := GetFaceProperties(validFace)
	if err != nil {
		t.Fatalf("could not get face properties: %v", err)
	}

	if properties.StyleName != "Book" {
		t.Errorf("expected StyleName to be \"Book\", got: %s", properties.StyleName)
	}

	if !properties.Monospace {
		t.Errorf("expected Monospace to be true, got: %t", properties.Monospace)
	}

	// invalid face test
	properties, err = GetFaceProperties(invalidFace)
	if err == nil {
		t.Fatalf("expected error with invalid face index (1234)")
	}
}

func writeRenderedCharacters(charsetName string, rendered RenderOutput, fontSize int, font string) {
	outDir := testutils.TmpPath(filepath.Join("rendered/bitmaps", charsetName))
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create directory %s: %v\n", outDir, err)
		return
	}

	for _, char := range rendered.Characters {
		if char.BitmapWidth <= 0 || char.BitmapHeight <= 0 || char.BitmapBuffer == nil {
			fmt.Printf("Skipping empty character U+%04X\n", char.Charcode)
			continue
		}

		unicodeStr := fmt.Sprintf("U+%04X", char.Charcode)
		fileName := fmt.Sprintf("%s_%s_%dpx.png", unicodeStr, font, fontSize)
		filePath := path.Join(outDir, fileName)

		img := image.NewGray(image.Rect(0, 0, char.BitmapWidth, char.BitmapHeight))
		img.Pix = char.BitmapBuffer

		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Printf("error writing file %s: %v\n", fileName, err)
			break
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			fmt.Printf("error exporting png into %s: %v\n", fileName, err)
			break
		}
	}
}

func TestRenderCharacters(t *testing.T) {
	supportedCharsets := []string{
		"alphanumeric",
		"ascii",
	}

	for _, charsetFile := range supportedCharsets {
		charsetPath := testutils.DataPath(filepath.Join("charsets", charsetFile+".txt"))
		charsetData, err := os.ReadFile(charsetPath)
		if err != nil {
			t.Fatalf("SETUP FAILURE: could not load charset file: %v", err)
		}

		charsetStr := string(charsetData)
		charset := []rune(charsetStr)

		rendered, error := Render(validFace, 16, charset)
		if error != nil {
			t.Fatalf("error during rendering characters: %v", error)
		}

		if len(rendered.Characters) == 0 {
			t.Fatalf("failed to render characters")
		} else if len(rendered.Characters[0].BitmapBuffer) != 0 {
			t.Errorf("whitespace bitmap should be empty")
		} else if rendered.Characters[0].Advance <= 0 {
			t.Errorf("whitespace advance expected to be larger than zero, got: %d", rendered.Characters[0].Advance)
		}

		for i := 1; i < len(rendered.Characters); i++ {
			if rendered.Characters[i].BitmapBuffer == nil {
				t.Errorf("bitmap must be non-empty for character %c", charset[i])
			}
		}

		writeRenderedCharacters(charsetFile, rendered, 16, "DejaVuSansMono")
	}

	unsupportedCharsets := []string{
		"braille",
	}

	for _, charsetFile := range unsupportedCharsets {
		charsetPath := testutils.DataPath(filepath.Join("charsets", charsetFile+".txt"))
		charsetData, err := os.ReadFile(charsetPath)
		if err != nil {
			t.Fatalf("SETUP FAILURE: could not load charset file: %v", err)
		}

		charsetStr := string(charsetData)
		charset := []rune(charsetStr)

		rendered, error := Render(validFace, 16, charset)
		if error != nil {
			t.Fatalf("error during rendering characters: %v", error)
		}

		if len(rendered.Characters) != 0 {
			t.Fatalf("there should be no rendered characters for unsupported charset, got: %d", len(rendered.Characters))
		}
		if len(rendered.Errors) != len(charset) {
			t.Errorf("for unsupported charset %d errors expected, got: %d", len(charset), len(rendered.Errors))
		}
	}
}
