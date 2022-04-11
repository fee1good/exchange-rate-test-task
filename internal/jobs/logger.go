package jobs

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type cronLogger struct {
	logger *zerolog.Logger
}

func newCronLogger(logger *zerolog.Logger) cron.Logger {
	return &cronLogger{logger: logger}
}

func (l *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info().Msgf(msg, keysAndValues...)
}

func (l *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error().Msgf(msg, keysAndValues...)
}
