package options

import "testing"

func TestHeightMatch_ValidString(t *testing.T) {
	chain := NewHeightChain(Terminal{})
	for _, s := range []string{"10px", "10l", "original", "th", "  10px     "} {
		if !chain.Match(s) {
			t.Errorf("expected %q to be a valid height", s)
		}
	}
}

func TestHeightMatch_InvalidString(t *testing.T) {
	chain := NewHeightChain(Terminal{})
	for _, s := range []string{"", "  ", "tw", "10 px", "01px", "0px", "-1px", "1.1px", "0l", "00l"} {
		if chain.Match(s) {
			t.Errorf("expected %q to be an invalid height", s)
		}
	}
}