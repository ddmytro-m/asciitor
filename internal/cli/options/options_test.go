package options

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v3"
)

func tempInput(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "in.png")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create temp input: %v", err)
	}
	return path
}

func parse(t *testing.T, args ...string) Values {
	t.Helper()

	var values Values
	cmd := &cli.Command{
		Arguments: []cli.Argument{&cli.StringArg{Name: "input"}},
		Flags:     Flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var err error
			values, err = Parse(cmd)
			return err
		},
	}

	if err := cmd.Run(context.Background(), append([]string{"asciitor"}, args...)); err != nil {
		t.Fatalf("unexpected error running command: %v", err)
	}
	return values
}

func TestParse_Defaults(t *testing.T) {
	values := parse(t, tempInput(t))
	defer values.Input.Close()
	defer values.Output.Close()

	if values.Input == nil {
		t.Error("expected a non-nil input reader")
	}
	if values.Output == nil {
		t.Error("expected a non-nil output writer")
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
	out := filepath.Join(t.TempDir(), "out.txt")
	values := parse(t,
		"-o", out,
		"-w", "100px",
		"-h", "50l",
		"--fill",
		"--inverse",
		tempInput(t),
	)
	defer values.Input.Close()
	defer values.Output.Close()

	if values.Input == nil {
		t.Error("expected a non-nil input reader")
	}
	if values.Output == nil {
		t.Error("expected a non-nil output writer")
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
