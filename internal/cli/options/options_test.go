package options

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"
)

func parse(t *testing.T, args ...string) Values {
	t.Helper()

	var values Values
	cmd := &cli.Command{
		Arguments: []cli.Argument{&cli.StringArg{Name: "input"}},
		Flags:     Flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			values = Parse(cmd)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), append([]string{"asciitor"}, args...)); err != nil {
		t.Fatalf("unexpected error running command: %v", err)
	}
	return values
}

func TestParse_Defaults(t *testing.T) {
	values := parse(t, "image.png")

	if values.Input != "image.png" {
		t.Errorf("expected input %q, got %q", "image.png", values.Input)
	}
	if values.Output != "-" {
		t.Errorf("expected default output %q, got %q", "-", values.Output)
	}
	if values.Width != "tw" {
		t.Errorf("expected default width %q, got %q", "tw", values.Width)
	}
	if values.Height != "th" {
		t.Errorf("expected default height %q, got %q", "th", values.Height)
	}
	if !values.KeepProportions {
		t.Error("expected KeepProportions to be true by default (fill is off)")
	}
	if values.Inverse {
		t.Error("expected Inverse to be false by default")
	}
}

func TestParse_Overrides(t *testing.T) {
	values := parse(t,
		"-o", "out.txt",
		"-w", "100px",
		"-h", "50l",
		"--fill",
		"--inverse",
		"in.png",
	)

	if values.Input != "in.png" {
		t.Errorf("expected input %q, got %q", "in.png", values.Input)
	}
	if values.Output != "out.txt" {
		t.Errorf("expected output %q, got %q", "out.txt", values.Output)
	}
	if values.Width != "100px" {
		t.Errorf("expected width %q, got %q", "100px", values.Width)
	}
	if values.Height != "50l" {
		t.Errorf("expected height %q, got %q", "50l", values.Height)
	}
	if values.KeepProportions {
		t.Error("expected KeepProportions to be false when --fill is set")
	}
	if !values.Inverse {
		t.Error("expected Inverse to be true when --inverse is set")
	}
}

func TestParse_DefaultsToStdin(t *testing.T) {
	values := parse(t)

	if values.Input != "" {
		t.Errorf("expected empty input (stdin) when no argument is given, got %q", values.Input)
	}
}

func TestParse_RejectsInvalidWidth(t *testing.T) {
	cmd := &cli.Command{
		Arguments: []cli.Argument{&cli.StringArg{Name: "input"}},
		Flags:     Flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	err := cmd.Run(context.Background(), []string{"asciitor", "-w", "bogus", "in.png"})
	if err == nil {
		t.Error("expected an error for an invalid width value")
	}
}
