package options

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/urfave/cli/v3"
)

var (
	rePx    = regexp.MustCompile(`^\d+px$`)
	reCols  = regexp.MustCompile(`^\d+[a-zA-Z]$`)
	reLines = regexp.MustCompile(`^\d+l$`)
)

func parseAmount(digits string) (int, bool) {
	if digits == "" || digits[0] == '0' {
		return 0, false
	}
	n, err := strconv.Atoi(digits)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

type Values struct {
	Input  io.ReadCloser
	Output io.WriteCloser

	Width  string
	Height string

	KeepProportions bool
	Inverse         bool
}

type Resolver[T, K any] interface {
	Resolve(T) (K, error)
}

type Matcher[T any] interface {
	Match(T) bool
}

func validate(m Matcher[string]) func(string) error {
	return func(s string) error {
		if !m.Match(s) {
			return fmt.Errorf("invalid value %q", s)
		}
		return nil
	}
}

var Arguments = []cli.Argument{
	&cli.StringArg{
		Name:      "input",
		UsageText: "input image file; omit or \"-\" to read from stdin",
	},
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:      "output",
		Aliases:   []string{"o"},
		Value:     "-",
		Usage:     "output file or a pipe (default to stdout)",
		Validator: validate(outputChain),
	},
	&cli.StringFlag{
		Name:      "width",
		Aliases:   []string{"w"},
		Value:     "tw",
		Usage:     "max output width: \"100px\", \"12M\" (12 letters M), \"original\" (image width), \"tw\" (terminal width - 1 line)",
		Validator: validate(NewWidthChain(Terminal{})),
	},
	&cli.StringFlag{
		Name:      "height",
		Aliases:   []string{"h"},
		Value:     "th",
		Usage:     "max output height: \"100px\", \"12l\" (12 lines), \"original\" (image height), \"th\" (terminal height)",
		Validator: validate(NewHeightChain(Terminal{})),
	},
	&cli.BoolWithInverseFlag{
		Name:  "fill",
		Value: false,
		Usage: "do not keep original proportions (may distort the image)",
	},
	&cli.BoolFlag{
		Name:  "inverse",
		Value: false,
		Usage: "inverse image colors",
	},
}

func Parse(cmd *cli.Command) (Values, error) {
	in := cmd.StringArg("input")
	if err := validate(inputChain)(in); err != nil {
		return Values{}, err
	}

	input, err := inputChain.Resolve(in)
	if err != nil {
		return Values{}, err
	}

	output, err := outputChain.Resolve(cmd.String("output"))
	if err != nil {
		input.Close()
		return Values{}, err
	}

	return Values{
		Input:  input,
		Output: output,

		Width:  cmd.String("width"),
		Height: cmd.String("height"),

		KeepProportions: !cmd.Bool("fill"),
		Inverse:         cmd.Bool("inverse"),
	}, nil
}
