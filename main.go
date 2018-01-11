/*
Command-line tool for creating and accessing backups.

Usage:

  $ kopia [<flags>] <subcommand> [<args> ...]

Use 'kopia help' to see more details.
*/
package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/kopia/kopia/cli"
	"github.com/kopia/kopia/repo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/alecthomas/kingpin.v2"

	_ "github.com/kopia/kopia/cli/filesystemcli"
	_ "github.com/kopia/kopia/cli/gcscli"
	_ "github.com/kopia/kopia/cli/webdavcli"
)

var (
	logFile  = cli.App().Flag("log-file", "log file name").String()
	logLevel = cli.App().Flag("log-level", "log level").Default("info").Enum("debug", "info", "warning", "error")
)

func initializeLogging(ctx *kingpin.ParseContext) error {
	switch *logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	if lfn := *logFile; lfn != "" {
		lf, err := os.Create(lfn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't create log file: %v", err)
			os.Exit(1)
		}

		log.Logger = log.Output(lf)
	} else {

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStderr()})
	}

	return nil
}

func main() {
	app := cli.App()
	app.Version(repo.BuildVersion + " build: " + repo.BuildInfo)
	app.PreAction(initializeLogging)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	return
}
