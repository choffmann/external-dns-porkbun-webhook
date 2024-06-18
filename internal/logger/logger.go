package logger

import (
	"errors"
	"io"
	"log/slog"
	"os"
)

type LogFormat string

const (
	Text LogFormat = "text"
	JSON LogFormat = "json"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func (l LogLevel) ToSLog() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func GetLogger(logFmt LogFormat) (*slog.Logger, error) {
	var w io.Writer = os.Stdout
	var options *slog.HandlerOptions

	var handler slog.Handler
	switch logFmt {
	case Text:
		handler = slog.NewTextHandler(w, options)
	case JSON:
		handler = slog.NewJSONHandler(w, options)
	default:
		return nil, errors.New("invalid logFormat!")
	}

	return slog.New(handler), nil
}
