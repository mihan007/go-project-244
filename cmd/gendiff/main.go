package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func GenDiff(ctx context.Context, cmd *cli.Command) error {
	format := cmd.String("format")
	if cmd.NArg() > 0 {
		filepath1 := cmd.Args().Get(0)
		filepath2 := cmd.Args().Get(1)
		result, err := code.GenDiff(filepath1, filepath2, format)
		if err != nil {
			return err
		}
		fmt.Println(result)
	} else {
		fmt.Println(cmd.Usage)
	}
	return nil
}

func main() {
	cmd := &cli.Command{
		Name:   "gendiff",
		Usage:  "Compares two configuration files and shows a difference.",
		Action: GenDiff,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Value:   "stylish",
				Usage:   "output format",
				Aliases: []string{"f"},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
