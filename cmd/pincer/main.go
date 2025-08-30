package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/milsim-tools/pincer/pkg/pincer"
	"github.com/urfave/cli/v2"
)

var Version = "dev"

func main() {
	app := cli.App{
		Name:  "pincer",
		Usage: "The milsim.tools platform CLI",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "The minimum level for logs to be emitted",
				EnvVars: []string{"HEARTH_LOG_LEVEL"},
				Aliases: []string{"l"},
				Value:   "info",
				Action: func(ctx *cli.Context, s string) error {
					switch s {
					case "debug":
					case "info":
					case "warn":
					case "error":
						return nil
					default:
						return fmt.Errorf("log-level must be one of debug, info, warn, error")
					}
					return nil
				},
			},
		},

		Commands: []*cli.Command{
			runCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
}

func runCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "run",
		Usage: "Starts the service(s) defined by configuration",
		Flags: pincer.Flags,
		Action: func(ctx *cli.Context) error {
			config := pincer.ConfigFromFlags(Version, ctx)

			var level slog.Level
			switch ctx.String("log-level") {
			case "debug":
				level = slog.LevelDebug
			case "info":
				level = slog.LevelInfo
			case "warn":
				level = slog.LevelWarn
			case "error":
				level = slog.LevelError
			}

			logger := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level: level,
				}),
			)

			h, err := pincer.New(logger, config)
			if err != nil {
				return err
			}

			return h.Run(ctx.Context)
		},
	}
	return cmd
}
