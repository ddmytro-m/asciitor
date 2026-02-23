package freetype

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"testing"

	"github.com/ddmytro-m/asciitor/internal/testutils"
)

var validFace = FaceParams{FaceIndex: 0, FontSize: 16}
var dummyFace = FaceParams{FaceIndex: 0, FontSize: 16}
var invalidFace = FaceParams{FaceIndex: 0, FontSize: -1}

func TestMain(m *testing.M) {
	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		fmt.Println("SETUP FAILURE: could not load font file:", err)
		os.Exit(1)
	}

	validFace.FontBuffer = fontData
	invalidFace.FontBuffer = fontData

	dummyPath := testutils.DataPath("fonts/dummy.ttf")
	dummyData, err := os.ReadFile(dummyPath)
	if err != nil {
		fmt.Println("SETUP FAILURE: could not load dummy font file:", err)
		os.Exit(1)
	}

	dummyFace.FontBuffer = dummyData

	m.Run()
}

func TestGetFaceProperties(t *testing.T) {
	// positive test
	properties, err := GetFaceProperties(validFace)
	if err != nil {
		t.Fatalf("could not get face properties: %v", err)
	}

	if properties.FamilyName != "DejaVu Sans Mono" {
		t.Errorf("expected FamilyName to be \"DejaVu Sans Mono\", got: %s", properties.FamilyName)
	}
	if properties.StyleName != "Book" {
		t.Errorf("expected StyleName to be \"Book\", got: %s", properties.StyleName)
	}

	if !properties.Monospace {
		t.Errorf("expected Monospace to be true, got: %t", properties.Monospace)
	}

	if properties.MaxCharacterWidth <= 0 {
		t.Errorf("expected MaxCharacterWidth to be bigger than 0, got: %d", properties.MaxCharacterWidth)
	}
	if properties.MaxCharacterHeight <= 0 {
		t.Errorf("expected MaxCharacterHeight to be bigger than 0, got: %d", properties.MaxCharacterHeight)
	}

	// dummy test
	properties, err = GetFaceProperties(dummyFace)
	if err == nil {
		t.Fatalf("expected error with invalid font file")
	}

	// invalid size test
	properties, err = GetFaceProperties(invalidFace)
	if err == nil {
		t.Fatalf("expected error with invalid font size (-1)")
	}
}

func writeRenderedCharacters(charset []rune, rendered []*RenderedCharacter, fontSize int, font string) {
	outDir := testutils.TmpPath("rendered/freetypeBitmaps")
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create directory %s: %v\n", outDir, err)
		return
	}

	for i := range len(charset) {
		renderedChar := rendered[i]
		if renderedChar == nil || renderedChar.BitmapBuffer == nil {
			continue
		}

		char := charset[i]
		unicodeStr := fmt.Sprintf("U+%04X", char)
		fileName := fmt.Sprintf("%s_%s_%dpx.png", unicodeStr, font, fontSize)
		filePath := path.Join(outDir, fileName)

		img := image.NewGray(image.Rect(0, 0, renderedChar.BitmapWidth, renderedChar.BitmapHeight))
		img.Pix = renderedChar.BitmapBuffer

		f, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
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
	faceParams := FaceParams{FaceIndex: 0, FontSize: 16}

	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatalf("SETUP FAILURE: could not load font file: %v", err)
	}
	faceParams.FontBuffer = fontData

	charsetPath := testutils.DataPath("charsets/alphanumeric.txt")
	charsetData, err := os.ReadFile(charsetPath)
	if err != nil {
		t.Fatalf("SETUP FAILURE: could not load font file: %v", err)
	}

	charsetStr := string(charsetData)
	charset := []rune(charsetStr)

	characters, error := RenderCharacters(faceParams, charset)
	if error != nil {
		t.Fatalf("error during rendering characters: %v", error)
	}

	if characters[0] == nil {
		t.Errorf("failed to render whitespace")
	} else if len(characters[0].BitmapBuffer) != 0 {
		t.Errorf("whitespace bitmap should be empty")
	} else if characters[0].Advance <= 0 {
		t.Errorf("whitespace advance expected to be larger than zero, got: %d", characters[0].Advance)
	}

	for i := 1; i < len(characters); i++ {
		if characters[i].BitmapBuffer == nil {
			t.Errorf("bitmap must be non-empty for character %c", charset[i])
		}
	}

	writeRenderedCharacters(charset, characters, 16, "DejaVuSansMono")
}
