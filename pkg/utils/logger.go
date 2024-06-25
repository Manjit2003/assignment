package utils

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func GetChildLogger(p string) *slog.Logger {
	return Logger.With("package", p)
}
