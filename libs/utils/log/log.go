package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

func SetupZeroLog() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
}

func SetupZeroLogWithCtx(ctx context.Context) (context.Context, zerolog.Logger) {
	logger := SetupZeroLog()

	return logger.WithContext(ctx), logger
}
