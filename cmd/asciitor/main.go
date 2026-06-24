package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

			in, err := options.OpenInput(values.Input)
			if err != nil {
				return err
			}
			defer in.Close()

			out, err := options.OpenOutput(values.Output)
			if err != nil {
				return err
			}
			defer out.Close()

			fmt.Printf("options: %+v\n", values)

			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
