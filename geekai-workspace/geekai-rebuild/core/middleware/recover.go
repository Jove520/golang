package middleware

import (
	"net/http"
	"runtime/debug"

	"geekai-rebuild/core/types"
	"geekai-rebuild/logger"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("handler panic: %v\n%s", r, string(debug.Stack()))

				c.JSON(http.StatusBadRequest, types.BizVo{
					Code:    types.Failed,
					Message: types.ErrorMsg,
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
