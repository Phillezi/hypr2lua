//go:build ignore

package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"

	"github.com/phillezi/gob"
)

func init() {
	logger := slog.New(gob.NewPrettyHandler(nil))
	slog.SetDefault(logger)
}

func main() {
	// takes in options
	gob.New(gob.WithDefaultTarget("all")).Add(
		"all",
		gob.Static(),
	).Add(
		"clean",
		gob.Clean(),
	).Add(
		"test",
		gob.NewCMD(func(ctx context.Context, cfg gob.Config) error {
			c := exec.CommandContext(
				ctx,
				"go",
				"test",
				"-v",
				"-race",
				"-timeout", "30s",
				"-vet=all",
				cfg.Selector,
			)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			return c.Run()
		}, "run tests"),
	).Add(
		"fix",
		gob.NewCMD(func(ctx context.Context, cfg gob.Config) error {
			c := exec.CommandContext(
				ctx,
				"go",
				"fix",
				cfg.Selector,
			)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			return c.Run()
		}, "fix code"),
	).Run()
}
