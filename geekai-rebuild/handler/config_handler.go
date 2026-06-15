package handler

import (
	"encoding/json"
	"errors"
	"strings"

	"geekai-rebuild/core"
	"geekai-rebuild/core/types"
	"geekai-rebuild/store/model"
	"geekai-rebuild/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	App          *core.AppServer
	DB           *gorm.DB
	SystemConfig *core.SystemConfigCache
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, systemConfig *core.SystemConfigCache) *ConfigHandler {
	return &ConfigHandler{
		App:          app,
		DB:           db,
		SystemConfig: systemConfig,
	}
}

func (h *ConfigHandler) RegisterRoutes() {
	api := h.App.Engine.Group("/api/config")
	api.GET("/get", h.GetConfig)
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if key == types.ConfigKeySystem {
		cfg := h.SystemConfig.Get()
		resp.SUCCESS(c, gin.H{
			"name":  key,
			"value": cfg.Base,
		})
		return
	}

	var config model.Config
	err := h.DB.Where("name = ?", key).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.ERROR(c, "配置不存在")
		return
	}
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var value any
	if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
		value = config.Value
	}

	resp.SUCCESS(c, gin.H{
		"name":  config.Name,
		"value": value,
	})
}
