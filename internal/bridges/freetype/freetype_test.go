package freetype

import (
	"os"
	"testing"

	"github.com/ddmytro-m/asciitor/internal/testutils"
)

func TestGetFaceProperties(t *testing.T) {
	validFontSize := 16
	invalidFontSize := -1

	fontPath := testutils.DataPath("fonts/DejaVuSansMono.ttf")
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatalf("SETUP FAILURE: could not load font file: %v", err)
	}

	dummyPath := testutils.DataPath("fonts/dummy.ttf")
	dummyData, err := os.ReadFile(dummyPath)
	if err != nil {
		t.Fatalf("SETUP FAILURE: could not load dummy font file: %v", err)
	}

	// positive test
	properties, err := GetFaceProperties(fontData, validFontSize)
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
	properties, err = GetFaceProperties(dummyData, validFontSize)
	if err == nil {
		t.Fatalf("expected error with invalid font file: %s", dummyPath)
	}

	// invalid size test
	properties, err = GetFaceProperties(fontData, invalidFontSize)
	if err == nil {
		t.Fatalf("expected error with invalid font size (%d)", invalidFontSize)
	}
}
