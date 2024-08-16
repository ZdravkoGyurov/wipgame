package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/app"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config: %s", err))
		os.Exit(1)
	}

	app, err := app.New(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to initialize app: %s", err))
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		slog.Error(fmt.Sprintf("failed to start app: %s", err))
		os.Exit(1)
	}
}
