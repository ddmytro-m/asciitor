package options

import (
	"path/filepath"
	"testing"
)

func TestOutputMatchers(t *testing.T) {
	if !(stdoutOutput{}).Match("-") {
		t.Error("stdout link must match \"-\"")
	}
	if (stdoutOutput{}).Match("result.txt") {
		t.Error("stdout link must not match a file path")
	}
	if !(fileOutput{}).Match("result.txt") {
		t.Error("file link must match a file path")
	}
	if (fileOutput{}).Match("-") {
		t.Error("file link must not match \"-\"")
	}
}

func TestOutputChain_Resolve(t *testing.T) {
	w, err := outputChain.Resolve("-")
	if err != nil {
		t.Fatalf("unexpected error resolving stdout: %v", err)
	}
	if w == nil {
		t.Error("expected a writer for stdout")
	}

	path := filepath.Join(t.TempDir(), "out.txt")
	f, err := outputChain.Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error resolving file: %v", err)
	}
	defer f.Close()
	if f == nil {
		t.Error("expected a writer for the file")
	}
}

type neverMatch struct{}

func (neverMatch) Match(string) bool { return false }

func TestValidate(t *testing.T) {
	v := validate(outputChain)
	if err := v("-"); err != nil {
		t.Errorf("unexpected error for valid output %q: %v", "-", err)
	}
	if err := v("result.txt"); err != nil {
		t.Errorf("unexpected error for valid output %q: %v", "result.txt", err)
	}

	if err := validate(neverMatch{})("anything"); err == nil {
		t.Error("expected an error when the matcher rejects the value")
	}
}