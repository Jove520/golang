package handler

import (
	"geekai-rebuild/core"
	"geekai-rebuild/store/model"
	"geekai-rebuild/store/vo"
	"geekai-rebuild/utils"
	"geekai-rebuild/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatModelHandler struct {
	App *core.AppServer
	DB  *gorm.DB
}

func NewChatModelHandler(app *core.AppServer, db *gorm.DB) *ChatModelHandler {
	return &ChatModelHandler{
		App: app,
		DB:  db,
	}
}

func (h *ChatModelHandler) RegisterRoutes() {
	api := h.App.Engine.Group("/api/chat")
	api.GET("/models", h.Models)
}

func (h *ChatModelHandler) Models(c *gin.Context) {
	var models []model.ChatModel

	if err := h.DB.
		Where("enabled = ?", true).
		Order("id ASC").
		Find(&models).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	items := make([]vo.ChatModel, 0, len(models))
	for _, item := range models {
		var modelVO vo.ChatModel
		_ = utils.CopyObject(&modelVO, item)
		items = append(items, modelVO)
	}

	resp.SUCCESS(c, items)
}
