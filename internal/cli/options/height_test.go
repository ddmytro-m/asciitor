package options

import "testing"

func TestHeightValidation_ValidString(t *testing.T) {
	err := validateHeight("10px")
	if err != nil {
		t.Error("unexpected error for valid pixel value")
	}

	err = validateHeight("10l")
	if err != nil {
		t.Error("unexpected error for valid lines value")
	}

	err = validateHeight("original")
	if err != nil {
		t.Error("unexpected error for original height")
	}

	err = validateHeight("th")
	if err != nil {
		t.Error("unexpected error for terminal height")
	}

	err = validateHeight("  10px     ")
	if err != nil {
		t.Error("whitespaces before and after the value must be ignored")
	}
}

func TestHeightValidation_InvalidString(t *testing.T) {
	err := validateHeight("")
	if err == nil {
		t.Error("expected an error for empty string")
	}

	err = validateHeight("  ")
	if err == nil {
		t.Error("expected an error for string consisting of only whitespaces")
	}

	err = validateHeight("tw")
	if err == nil {
		t.Error("expected an error for \"tw\"")
	}

	err = validateHeight("10 px")
	if err == nil {
		t.Error("whitespaces between value and unit are not allowed")
	}

	err = validateHeight("01px")
	if err == nil {
		t.Error("numbers cannot start with zero")
	}

	err = validateHeight("0px")
	if err == nil {
		t.Error("pixel value must be a positive integer")
	}

	err = validateHeight("-1px")
	if err == nil {
		t.Error("pixel value cannot be negative")
	}

	err = validateHeight("1.1px")
	if err == nil {
		t.Error("pixel value must be a positive integer")
	}

	err = validateHeight("0l")
	if err == nil {
		t.Error("lines count must be a positive integer")
	}

	err = validateHeight("00l")
	if err == nil {
		t.Error("zero lines count must be rejected")
	}
}
