package handler

import (
	"context"
	"fmt"
	"geekai-rebuild/core"
	"geekai-rebuild/utils"
	"geekai-rebuild/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CaptchaHandler struct {
	App   *core.AppServer
	Redis *redis.Client
}

func NewCaptchaHandler(app *core.AppServer, redisClient *redis.Client) *CaptchaHandler {
	return &CaptchaHandler{
		App:   app,
		Redis: redisClient,
	}
}

func (h *CaptchaHandler) RegisterRoutes() {
	api := h.App.Engine.Group("/api/captcha")
	api.GET("", h.Generate)
}

func (h *CaptchaHandler) Generate(c *gin.Context) {
	captchaId, err := utils.GenSalt(16)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	code, err := utils.GenNumericCode(6)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	key := fmt.Sprintf("captcha/%s", captchaId)
	if err := h.Redis.Set(context.Background(), key, code, 5*time.Minute).Err(); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"id":   captchaId,
		"code": code,
	})
}
