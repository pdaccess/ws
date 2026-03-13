package app

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	DatabaseConfig  DatabaseConfig
	PDAwsConfig     PDAwsConfig
	AuthwsConfig    AuthwsConfig
	CManangerConfig CManangerConfig

	HttpListenAddr     string `env:"HTTP_LISTEN_ADDR" envDefault:":8080"`
	SshListenAddr      string `env:"SSH_LISTEN_ADDR" envDefault:":2222"`
	TerminalListenAddr string `env:"TERMINAL_LISTEN_ADDR" envDefault:":8081"`
	GuacdAddr          string `env:"GUACD_ADDR" envDefault:"proxy:8081"`

	ProxyProtocol bool `env:"PROXY_PROTOCOL" envDefault:"false"`

	ServerKeyPath   string `env:"SSH_SERVER_KEY_PATH" envDefault:"/etc/security/ssh-host-key"`
	ConnectionLimit int    `env:"SSH_CONNECTION_LIMIT" envDefault:"100"`
	ConnectionRate  int    `env:"SSH_CONNECTION_RATE" envDefault:"3"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"INFO"`
}

type AuthwsConfig struct {
	Url string `env:"AUTHWS_URL" envDefault:"iauth:55051"`
}

type PDAwsConfig struct {
	Url string `env:"PDAWS_URL" envDefault:"http://pdaws:8080"`
}

type CManangerConfig struct {
	Url string `env:"CMANAGER_URL" envDefault:"cmanager:50051"`
}

type DatabaseConfig struct {
	Url string `env:"DATABASE_URL" envDefault:"postgresql://pda:pda@pda:5432/pda?sslmode=disable"`
}

func ParseConfig() (*ServerConfig, error) {
	cfg := ServerConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config parse: %w", err)
	}

	log.Info().
		Interface("config", cfg).
		Msgf("dump current config")

	return &cfg, nil
}

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
