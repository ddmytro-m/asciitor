package options

import (
	"reflect"
	"testing"

	"github.com/ddmytro-m/asciitor/sizing"
)

func TestOutputWidth(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		term    Terminal
		want    sizing.OutputWidth
		wantErr bool
	}{
		{name: "original", in: "original", want: sizing.WidthAuto{}},
		{name: "terminal width", in: "tw", term: Terminal{Cols: 80, Ref: 'M'}, want: sizing.WidthCharacters{Character: 'M', Amount: 80}},
		{name: "terminal width unknown falls back to auto", in: "tw", term: Terminal{Cols: 0, Ref: 'M'}, want: sizing.WidthAuto{}},
		{name: "terminal width negative falls back to auto", in: "tw", term: Terminal{Cols: -1}, want: sizing.WidthAuto{}},
		{name: "pixels", in: "100px", want: sizing.WidthPixels{Pixels: 100}},
		{name: "characters", in: "12M", want: sizing.WidthCharacters{Character: 'M', Amount: 12}},
		{name: "single character", in: "5a", want: sizing.WidthCharacters{Character: 'a', Amount: 5}},
		{name: "empty", in: "", wantErr: true},
		{name: "unknown", in: "bogus", wantErr: true},
		{name: "lines suffix is not a width", in: "50l", want: sizing.WidthCharacters{Character: 'l', Amount: 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWidthChain(tt.term).Resolve(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got %#v", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outputWidth(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestOutputHeight(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		term    Terminal
		want    sizing.OutputHeight
		wantErr bool
	}{
		{name: "original", in: "original", want: sizing.HeightAuto{}},
		{name: "terminal height reserves one row", in: "th", term: Terminal{Rows: 24}, want: sizing.HeightLines{Amount: 23}},
		{name: "terminal height of two yields one line", in: "th", term: Terminal{Rows: 2}, want: sizing.HeightLines{Amount: 1}},
		{name: "terminal height of one falls back to auto", in: "th", term: Terminal{Rows: 1}, want: sizing.HeightAuto{}},
		{name: "terminal height unknown falls back to auto", in: "th", term: Terminal{Rows: 0}, want: sizing.HeightAuto{}},
		{name: "pixels", in: "100px", want: sizing.HeightPixels{Pixels: 100}},
		{name: "lines", in: "50l", want: sizing.HeightLines{Amount: 50}},
		{name: "empty", in: "", wantErr: true},
		{name: "unknown", in: "qwerty", wantErr: true},
		{name: "character suffix is not a height", in: "12M", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHeightChain(tt.term).Resolve(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got %#v", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outputHeight(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestOutputSize(t *testing.T) {
	v := Values{Width: "100px", Height: "50l"}
	got, err := v.OutputSize(Terminal{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := sizing.OutputSize{Width: sizing.WidthPixels{Pixels: 100}, Height: sizing.HeightLines{Amount: 50}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OutputSize() = %#v, want %#v", got, want)
	}
}

func TestOutputSize_TrimsWhitespace(t *testing.T) {
	v := Values{Width: "  100px ", Height: "\t50l\n"}
	got, err := v.OutputSize(Terminal{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := sizing.OutputSize{Width: sizing.WidthPixels{Pixels: 100}, Height: sizing.HeightLines{Amount: 50}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OutputSize() = %#v, want %#v", got, want)
	}
}

func TestOutputSize_PropagatesErrors(t *testing.T) {
	if _, err := (Values{Width: "bogus", Height: "50l"}).OutputSize(Terminal{}); err == nil {
		t.Error("expected error from invalid width")
	}
	if _, err := (Values{Width: "100px", Height: "bogus"}).OutputSize(Terminal{}); err == nil {
		t.Error("expected error from invalid height")
	}
}