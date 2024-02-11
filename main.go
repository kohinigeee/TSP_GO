package main

import (
	"log/slog"
	"tspgo/cmd"
	mylogger "tspgo/logger"
)

func main() {
	// アルゴリズムの細かなログの表示レベル
	mylogger.SetAlgoLevel(slog.LevelError)

	cmd.MainWithGreeding()
	// cmd.MainWithPprof()
}
