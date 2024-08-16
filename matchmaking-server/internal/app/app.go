package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/matchmaking"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/storage/redis"
)

type GlobalContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func NewGlobalContext() GlobalContext {
	ctx, cancel := context.WithCancel(context.Background())

	return GlobalContext{
		Context: ctx,
		Cancel:  cancel,
	}
}

type Application struct {
	appContext        GlobalContext
	config            config.Config
	redisClient       redis.Client
	matchmakingServer matchmaking.Server
}

func New(cfg config.Config) (*Application, error) {
	globalContext := NewGlobalContext()

	redisClient := redis.NewClient(cfg.Redis)

	matchmakingServer := matchmaking.NewServer(cfg.Matchmaking, redisClient)

	return &Application{
		appContext:        globalContext,
		config:            cfg,
		redisClient:       redisClient,
		matchmakingServer: matchmakingServer,
	}, nil
}

func (a *Application) Start() error {
	slog.Info("starting application...")
	a.setupSignalNotifier()

	go func() {
		a.matchmakingServer.Start(a.appContext.Context)
	}()

	slog.Info("application started")

	<-a.appContext.Context.Done()

	a.stopMatchmakingServer()
	a.closeRedisClient()
	slog.Info("application stopped")

	return nil
}

func (a *Application) Stop() {
	a.appContext.Cancel()
}

func (a *Application) stopMatchmakingServer() {
	a.matchmakingServer.Stop()
	slog.Info("matchmaking server stopped")
}

func (a *Application) closeRedisClient() {
	a.redisClient.Close()
	slog.Info("redis client closed")
}

func (a *Application) setupSignalNotifier() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		slog.Info("stopping application...")
		a.appContext.Cancel()
	}()
}
