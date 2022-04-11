package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func New(logLevel string) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse zerolog level: %w", err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zeroLogger := zerolog.New(os.Stdout)
	zeroLogger.Level(level)

	return &zeroLogger, nil
}
