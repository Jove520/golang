package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"geekai-rebuild/core"
	"geekai-rebuild/core/types"
	"geekai-rebuild/handler"
	"geekai-rebuild/logger"
	"geekai-rebuild/service"
	"geekai-rebuild/store"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			func() (*types.AppConfig, error) {
				return core.LoadConfig("config.toml")
			},
			core.NewServer,
			core.NewSystemConfigCache,
			handler.NewBaseHandler,
			handler.NewConfigHandler,
			handler.NewUserHandler,
			handler.NewCaptchaHandler,
			handler.NewChatModelHandler,
			logger.GetLogger,
			store.NewGormConfig,
			store.NewMysql,
			store.NewRedisClient,
			service.NewApiKeyService,
		),
		fx.Invoke(store.AutoMigrateAndSeed),
		fx.Invoke(store.CheckRedis),
		fx.Invoke(func(h *handler.BaseHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(h *handler.ConfigHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(h *handler.UserHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(h *handler.CaptchaHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(h *handler.ChatModelHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(startServer),
	)

	go func() {
		if err := app.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}

func startServer(lifecycle fx.Lifecycle, server *core.AppServer, log *zap.SugaredLogger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalw("server stopped", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down server")
			return server.Shutdown(ctx)
		},
	})
}
