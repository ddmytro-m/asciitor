package options

import (
	"os"
	"path/filepath"
	"testing"
)

type fakeFontRepository struct {
	fonts map[string][]byte
}

func (f fakeFontRepository) Has(val string) bool {
	_, ok := f.fonts[val]
	return ok
}

func (f fakeFontRepository) Get(val string) ([]byte, error) {
	data, ok := f.fonts[val]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

func TestFontMatchers(t *testing.T) {
	bundled := fontBundled{repository: fakeFontRepository{fonts: map[string][]byte{"mono": {0x01}}}}
	if !bundled.Match("mono") {
		t.Error("bundled link must match a known font name")
	}
	if bundled.Match("unknown") {
		t.Error("bundled link must not match an unknown font name")
	}

	if !(fontFile{}).Match("font.ttf") {
		t.Error("file link must match a file path")
	}
	if (fontFile{}).Match("") {
		t.Error("file link must not match an empty argument")
	}
}

func TestFontBundled_Resolve(t *testing.T) {
	bundled := fontBundled{repository: fakeFontRepository{fonts: map[string][]byte{"mono": {0x01, 0x02, 0x03}}}}
	got, err := bundled.Resolve("mono")
	if err != nil {
		t.Fatalf("unexpected error resolving bundled font: %v", err)
	}
	if string(got) != string([]byte{0x01, 0x02, 0x03}) {
		t.Errorf("expected %v, got %v", []byte{0x01, 0x02, 0x03}, got)
	}
}

func TestFontFile_Resolve(t *testing.T) {
	path := filepath.Join(t.TempDir(), "font.ttf")
	want := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	if err := os.WriteFile(path, want, 0o644); err != nil {
		t.Fatalf("failed to create temp font: %v", err)
	}

	got, err := (fontFile{}).Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFontFile_MissingFile(t *testing.T) {
	if _, err := (fontFile{}).Resolve(filepath.Join(t.TempDir(), "missing.ttf")); err == nil {
		t.Error("expected an error for a missing font file")
	}
}

func TestFontChain_ResolvesBundled(t *testing.T) {
	got, err := fontChain.Resolve("dejavusansmono")
	if err != nil {
		t.Fatalf("unexpected error resolving bundled font: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected non-empty data for the bundled font")
	}
}

func TestFontChain_ResolvesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "font.ttf")
	want := []byte{0x00, 0x01, 0x00, 0x00}
	if err := os.WriteFile(path, want, 0o644); err != nil {
		t.Fatalf("failed to create temp font: %v", err)
	}

	got, err := fontChain.Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file font: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFontChain_Unresolvable(t *testing.T) {
	if _, err := fontChain.Resolve(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("expected an error for a value that is neither a bundled font nor a file")
	}
}