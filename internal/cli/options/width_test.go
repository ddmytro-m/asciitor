package options

import "testing"

func TestWidthMatch_ValidString(t *testing.T) {
	chain := NewWidthChain(Terminal{})
	for _, s := range []string{"10px", "12M", "original", "tw", "  10px     "} {
		if !chain.Match(s) {
			t.Errorf("expected %q to be a valid width", s)
		}
	}
}

func TestWidthMatch_InvalidString(t *testing.T) {
	chain := NewWidthChain(Terminal{})
	for _, s := range []string{"", "  ", "th", "12 M", "0M", "00M", "0px", "00px", "-1px", "1.1px", "abc"} {
		if chain.Match(s) {
			t.Errorf("expected %q to be an invalid width", s)
		}
	}
}