package options

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInputMatchers(t *testing.T) {
	if !(stdinInput{}).Match("-") {
		t.Error("stdin link must match \"-\"")
	}
	if !(stdinInput{}).Match("") {
		t.Error("stdin link must match an empty argument")
	}
	if (stdinInput{}).Match("image.png") {
		t.Error("stdin link must not match a file path")
	}
	if !(fileInput{}).Match("image.png") {
		t.Error("file link must match a file path")
	}
	if (fileInput{}).Match("-") {
		t.Error("file link must not match \"-\"")
	}
}

func TestInputChain_ResolvesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in.png")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create temp input: %v", err)
	}

	r, err := inputChain.Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file: %v", err)
	}
	defer r.Close()
	if r == nil {
		t.Error("expected a reader for the file")
	}
}

func TestInputChain_MissingFile(t *testing.T) {
	if _, err := inputChain.Resolve(filepath.Join(t.TempDir(), "missing.png")); err == nil {
		t.Error("expected an error for a missing input file")
	}
}