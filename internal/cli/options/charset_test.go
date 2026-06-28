package options

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCharsetMatchers(t *testing.T) {
	preset := charsetPreset{presets: map[string]string{"ascii": " .:#"}}
	if !preset.Match("ascii") {
		t.Error("preset link must match a known preset name")
	}
	if preset.Match("unknown") {
		t.Error("preset link must not match an unknown preset name")
	}

	if !(charsetFile{}).Match("charset.txt") {
		t.Error("file link must match a file path")
	}
	if (charsetFile{}).Match("") {
		t.Error("file link must not match an empty argument")
	}
}

func TestCharsetPreset_Resolve(t *testing.T) {
	preset := charsetPreset{presets: map[string]string{"ascii": " .:#"}}
	got, err := preset.Resolve("ascii")
	if err != nil {
		t.Fatalf("unexpected error resolving preset: %v", err)
	}
	if string(got) != " .:#" {
		t.Errorf("expected %q, got %q", " .:#", string(got))
	}
}

func TestCharsetFile_Resolve(t *testing.T) {
	path := filepath.Join(t.TempDir(), "charset.txt")
	if err := os.WriteFile(path, []byte("@%#*+=-:. "), 0o644); err != nil {
		t.Fatalf("failed to create temp charset: %v", err)
	}

	got, err := (charsetFile{}).Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file: %v", err)
	}
	if string(got) != "@%#*+=-:. " {
		t.Errorf("expected %q, got %q", "@%#*+=-:. ", string(got))
	}
}

func TestCharsetFile_MissingFile(t *testing.T) {
	if _, err := (charsetFile{}).Resolve(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Error("expected an error for a missing charset file")
	}
}

func TestCharsetChain_ResolvesBuiltin(t *testing.T) {
	got, err := charsetChain.Resolve("ascii")
	if err != nil {
		t.Fatalf("unexpected error resolving builtin charset: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected a non-empty charset for the ascii builtin")
	}
}

func TestCharsetChain_ResolvesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "charset.txt")
	if err := os.WriteFile(path, []byte(" .oO@"), 0o644); err != nil {
		t.Fatalf("failed to create temp charset: %v", err)
	}

	got, err := charsetChain.Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file charset: %v", err)
	}
	if string(got) != " .oO@" {
		t.Errorf("expected %q, got %q", " .oO@", string(got))
	}
}

func TestCharsetChain_Unresolvable(t *testing.T) {
	if _, err := charsetChain.Resolve(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("expected an error for a value that is neither a preset nor a file")
	}
}