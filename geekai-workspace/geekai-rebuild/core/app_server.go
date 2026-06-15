package core

import (
	"context"
	"net/http"

	"geekai-rebuild/core/middleware"
	"geekai-rebuild/core/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppServer struct {
	Config *types.AppConfig
	Engine *gin.Engine
	Log    *zap.SugaredLogger
	server *http.Server
}

func NewServer(config *types.AppConfig, log *zap.SugaredLogger) *AppServer {
	gin.SetMode(gin.ReleaseMode)

	server := &AppServer{
		Config: config,
		Engine: gin.New(),
		Log:    log,
	}

	server.registerMiddlewares()

	return server
}

func (s *AppServer) Run() error {
	s.server = &http.Server{
		Addr:    s.Config.Listen,
		Handler: s.Engine,
	}

	s.Log.Infof("server listening on http://%s", s.Config.Listen)
	return s.server.ListenAndServe()
}

func (s *AppServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *AppServer) registerMiddlewares() {
	s.Engine.Use(middleware.TrimSpaceMiddleware())
	s.Engine.Use(middleware.RecoverMiddleware())
}
