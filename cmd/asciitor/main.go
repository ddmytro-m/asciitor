package main

import (
	"context"
	"log"
	"os"

	"github.com/ddmytro-m/asciitor/internal/cli/app"
	"github.com/ddmytro-m/asciitor/internal/cli/options"
	"github.com/urfave/cli/v3"
)

func main() {
	var cmd = &cli.Command{
		Name:      "asciitor",
		Usage:     "convert an image to ASCII art",
		ArgsUsage: "[INPUT]",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "input",
				UsageText: "input image file; omit or \"-\" to read from stdin",
			},
		},
		Flags: options.Flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			values := options.Parse(cmd)
			return app.Run(ctx, values)
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
