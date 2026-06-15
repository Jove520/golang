package handler

import (
	"geekai-rebuild/core"
	"geekai-rebuild/core/types"
	"geekai-rebuild/utils/resp"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	App *core.AppServer
}

func NewBaseHandler(app *core.AppServer) *BaseHandler {
	return &BaseHandler{App: app}
}

func (h *BaseHandler) RegisterRoutes() {
	api := h.App.Engine.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		resp.SUCCESS(c, "pong")
	})

	api.GET("/panic", func(c *gin.Context) {
		panic("manual panic test")
	})

	api.POST("/echo", func(c *gin.Context) {
		var body any

		if err := c.ShouldBindJSON(&body); err != nil {
			resp.ERROR(c, types.InvalidArgs)
			return
		}

		resp.SUCCESS(c, body)
	})
}
