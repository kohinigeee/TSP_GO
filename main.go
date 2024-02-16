package main

import (
	"log/slog"
	"tspgo/cmd"
	mylogger "tspgo/logger"
)

func main() {
	// アルゴリズムの細かなログの表示レベル
	mylogger.SetAlgoLevel(slog.LevelError)
	mylogger.SetLevel(slog.LevelInfo)

	// cmd.MainWithLocalSearch()
	// cmd.MainWithGreeding()
	// cmd.MainWithPprof()

	cmd.MainCtest()
}
