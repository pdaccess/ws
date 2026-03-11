package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pdaccess/ws/cmd/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Commit, BuildTime, BuildEnv string

func setupLog(debug, console bool) {
	zerolog.TimeFieldFormat = "2006-01-10,10:01:02"

	zerolog.TimestampFunc = func() time.Time {
		return time.Now()
	}
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	if console {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp().Logger().Level(level)
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
	}
}

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	console := flag.Bool("console", false, "sets log level to debug")

	flag.Parse()
	setupLog(*debug, *console)

	log.Info().
		Str("Commnit", Commit).
		Str("BuildTime", BuildTime).
		Str("nBuildEnv", BuildEnv).
		Msg("WS Starting")

	log.Info().Msgf("log level is: %s", zerolog.GlobalLevel())

	config, err := app.ParseConfig()

	if err != nil {
		log.Err(err).Msg("read config")
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := app.StartServer(config, signalChan); err != nil {
		log.Err(err).Msg("start server")
		os.Exit(1)
	}
}
