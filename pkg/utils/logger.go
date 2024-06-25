package utils

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

func GetChildLogger(p string) *slog.Logger {
	return Logger.With("package", p)
}
