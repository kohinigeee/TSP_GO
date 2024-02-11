package mylogger

import (
	"log/slog"
	"os"

	"github.com/kohinigeee/mylog/clog"
)

//アルゴリズム中の細かなWarnログを表示するかどうかを設定する

var (
	algoLogger   *slog.Logger
	algoLogLevel *slog.LevelVar
)

func init() {
	algoLogLevel = new(slog.LevelVar)
	algoLogLevel.Set(slog.LevelDebug)

	handlerOption := &slog.HandlerOptions{
		Level: algoLogLevel,
	}

	handler, err := clog.NewCustomTextHandler(os.Stdout,
		clog.WithHandlerOption(handlerOption))

	if err != nil {
		panic(err)
	}

	algoLogger = slog.New(handler)
}

func SetAlgoLevel(level slog.Level) {
	algoLogLevel.Set(level)
}

func AL() *slog.Logger {
	return algoLogger
}
