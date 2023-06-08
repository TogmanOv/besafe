package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(out io.Writer, pretty bool) zerolog.Logger {
	logger := zerolog.New(os.Stdout)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out, TimeFormat: "15:04:05"})
	}

	return logger
}
