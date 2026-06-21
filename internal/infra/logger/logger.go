package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// ANSI Color Codes for pretty local development logging
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorCyan   = "\033[36m"
)

// PrettyHandler is a custom slog.Handler that formats logs nicely for console reading
type PrettyHandler struct {
	slog.Handler
	out io.Writer
}

func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		Handler: slog.NewTextHandler(out, opts),
		out:     out,
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelStr := r.Level.String()
	var levelColor string

	switch r.Level {
	case slog.LevelDebug:
		levelColor = colorGray
	case slog.LevelInfo:
		levelColor = colorBlue
	case slog.LevelWarn:
		levelColor = colorYellow
	case slog.LevelError:
		levelColor = colorRed
	default:
		levelColor = colorReset
	}

	// Format timestamp
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")

	// Print timestamp, colored level, and message
	_, _ = h.out.Write([]byte(
		colorGray + "[" + timeStr + "] " + colorReset +
			levelColor + "[" + levelStr + "] " + colorReset +
			r.Message,
	))

	// Print attributes
	r.Attrs(func(a slog.Attr) bool {
		valStr := a.Value.String()
		// Highlight some standard keys
		keyColor := colorCyan
		if strings.Contains(a.Key, "err") || strings.Contains(a.Key, "error") {
			keyColor = colorRed
		}
		_, _ = h.out.Write([]byte(" " + keyColor + a.Key + colorReset + "=" + valStr))
		return true
	})

	_, _ = h.out.Write([]byte("\n"))
	return nil
}

// SetupLogger initializes the global structured logger
func SetupLogger(env, levelStr string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(env) == "local" {
		handler = NewPrettyHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
