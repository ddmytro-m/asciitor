package options

import "testing"

func TestOutputValidation_ValidString(t *testing.T) {
	err := validateOutput("result.txt")
	if err != nil {
		t.Error("unexpected error for a valid output path")
	}

	err = validateOutput("-")
	if err != nil {
		t.Error("unexpected error for stdout output")
	}

	err = validateOutput("  result.txt  ")
	if err != nil {
		t.Error("whitespaces before and after the value must be ignored")
	}
}

func TestOutputValidation_InvalidString(t *testing.T) {
	err := validateOutput("")
	if err == nil {
		t.Error("expected an error for empty string")
	}

	err = validateOutput("   ")
	if err == nil {
		t.Error("expected an error for string consisting of only whitespaces")
	}
}
