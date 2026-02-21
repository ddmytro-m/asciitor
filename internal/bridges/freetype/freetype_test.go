package freetype

import (
	"fmt"
	"os"
	"testing"

	"github.com/ddmytro-m/asciitor/internal/testutils"
)

func TestGetFaceProperties(t *testing.T) {
	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	fontSize := 16

	data, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatalf("SETUP FAILURE: could not load font file: %v", err)
	}

	fmt.Println(len(data))

	properties, err := GetFaceProperties(data, fontSize)
	if err != nil {
		t.Fatalf("could not get face properties: %s", err)
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
}
