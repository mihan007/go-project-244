package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func GenDiff(ctx context.Context, cmd *cli.Command) error {
	fmt.Println(cmd.Usage)
	return nil
}

func main() {
	cmd := &cli.Command{
		Name:   "gendiff",
		Usage:  "Compares two configuration files and shows a difference.",
		Action: GenDiff,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
