package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
)

var logger *slog.Logger

type adjustedHandler struct {
    slog.Handler
}

func (h adjustedHandler) Handle(ctx context.Context, r slog.Record) error {
    pcs := make([]uintptr, 1)
    runtime.Callers(4, pcs) 
    fs := runtime.CallersFrames(pcs)
    frame, _ := fs.Next()
    r.PC = frame.PC

    return h.Handler.Handle(ctx, r)
}


func Init(env string) {
	var handler slog.Handler

	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	}

	logger = slog.New(adjustedHandler{ Handler: handler})
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}
