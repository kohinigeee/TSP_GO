package mylogger

import (
	"log/slog"
	"os"

	"github.com/kohinigeee/mylog/clog"
)

var (
	logLevel *slog.LevelVar
	logger   *slog.Logger
)

func L() *slog.Logger {
	return logger
}

func SetLevel(level slog.Level) {
	logLevel.Set(level)
}

func init() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	handlerOption := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler, err := clog.NewCustomTextHandler(os.Stdout,
		clog.WithHandlerOption(handlerOption))

	if err != nil {
		panic(err)
	}

	logger = slog.New(handler)
}
