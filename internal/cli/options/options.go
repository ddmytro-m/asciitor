package options

import (
	"regexp"

	"github.com/urfave/cli/v3"
)

var (
	rePx    = regexp.MustCompile(`^\d+px$`)
	reCols  = regexp.MustCompile(`^\d+[a-zA-Z]$`)
	reLines = regexp.MustCompile(`^\d+l$`)
)

type Values struct {
	Input  string
	Output string

	Width  string
	Height string

	KeepProportions bool
	Inverse         bool
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:      "output",
		Aliases:   []string{"o"},
		Value:     "-",
		Usage:     "output file or a pipe (default to stdout)",
		Validator: validateOutput,
	},
	&cli.StringFlag{
		Name:      "width",
		Aliases:   []string{"w"},
		Value:     "tw",
		Usage:     "max output width: \"100px\", \"12M\" (12 letters M), \"original\" (image width), \"tw\" (terminal width - 1 line)",
		Validator: validateWidth,
	},
	&cli.StringFlag{
		Name:      "height",
		Aliases:   []string{"h"},
		Value:     "th",
		Usage:     "max output height: \"100px\", \"12l\" (12 lines), \"original\" (image height), \"th\" (terminal height)",
		Validator: validateHeight,
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

func Parse(cmd *cli.Command) Values {
	return Values{
		Input:  cmd.StringArg("input"),
		Output: cmd.String("output"),

		Width:  cmd.String("width"),
		Height: cmd.String("height"),

		KeepProportions: !cmd.Bool("fill"),
		Inverse:         cmd.Bool("inverse"),
	}
}
