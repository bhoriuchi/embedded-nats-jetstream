package logging

import (
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/rs/zerolog"
)

const (
	natsSkipFrame = 0
)

func NewNATSLogger(logg zerolog.Logger) natsserver.Logger {
	return &natsLogger{
		l: logg,
	}
}

type natsLogger struct {
	l zerolog.Logger
}

// Log a notice statement
func (l *natsLogger) Noticef(format string, v ...interface{}) {
	l.l.Info().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}

// Log a warning statement
func (l *natsLogger) Warnf(format string, v ...interface{}) {
	l.l.Warn().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}

// Log a fatal error
func (l *natsLogger) Fatalf(format string, v ...interface{}) {
	l.l.Fatal().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}

// Log an error
func (l *natsLogger) Errorf(format string, v ...interface{}) {
	l.l.Error().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}

// Log a debug statement
func (l *natsLogger) Debugf(format string, v ...interface{}) {
	l.l.Debug().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}

// Log a trace statement
func (l *natsLogger) Tracef(format string, v ...interface{}) {
	l.l.Trace().CallerSkipFrame(natsSkipFrame).Msgf(format, v...)
}
