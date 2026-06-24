package options

import "testing"

func TestWidthValidation_ValidString(t *testing.T) {
	err := validateWidth("10px")
	if err != nil {
		t.Error("unexpected error for valid pixel value")
	}

	err = validateWidth("12M")
	if err != nil {
		t.Error("unexpected error for valid characters value")
	}

	err = validateWidth("original")
	if err != nil {
		t.Error("unexpected error for original width")
	}

	err = validateWidth("tw")
	if err != nil {
		t.Error("unexpected error for terminal width")
	}

	err = validateWidth("  10px     ")
	if err != nil {
		t.Error("whitespaces before and after the value must be ignored")
	}
}

func TestWidthValidation_InvalidString(t *testing.T) {
	err := validateWidth("")
	if err == nil {
		t.Error("expected an error for empty string")
	}

	err = validateWidth("  ")
	if err == nil {
		t.Error("expected an error for string consisting of only whitespaces")
	}

	err = validateWidth("th")
	if err == nil {
		t.Error("expected an error for \"th\"")
	}

	err = validateWidth("12 M")
	if err == nil {
		t.Error("whitespaces between value and unit are not allowed")
	}

	err = validateWidth("0M")
	if err == nil {
		t.Error("characters count must be a positive integer")
	}

	err = validateWidth("00M")
	if err == nil {
		t.Error("zero characters count must be rejected")
	}

	err = validateWidth("0px")
	if err == nil {
		t.Error("pixel value must be a positive integer")
	}

	err = validateWidth("00px")
	if err == nil {
		t.Error("zero pixel value must be rejected")
	}

	err = validateWidth("-1px")
	if err == nil {
		t.Error("pixel value cannot be negative")
	}

	err = validateWidth("1.1px")
	if err == nil {
		t.Error("pixel value must be a positive integer")
	}

	err = validateWidth("abc")
	if err == nil {
		t.Error("expected an error for an unknown width value")
	}
}
