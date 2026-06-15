package handler

import (
	"geekai-rebuild/core"
	"geekai-rebuild/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ChatStreamHandler struct {
	App           *core.AppServer
	DB            *gorm.DB
	Redis         *redis.Client
	ApiKeyService *service.ApiKeyService
}
