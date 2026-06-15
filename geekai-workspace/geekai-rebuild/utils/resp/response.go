package resp

import (
	"net/http"

	"geekai-rebuild/core/types"

	"github.com/gin-gonic/gin"
)

func SUCCESS(c *gin.Context, data ...any) {
	vo := types.BizVo{Code: types.Success}

	if len(data) > 0 {
		vo.Data = data[0]
	}

	c.JSON(http.StatusOK, vo)
}

func ERROR(c *gin.Context, message ...string) {
	vo := types.BizVo{Code: types.Failed}

	if len(message) > 0 {
		vo.Message = message[0]
	}

	c.JSON(http.StatusBadRequest, vo)
}

func NotAuth(c *gin.Context, message ...string) {
	vo := types.BizVo{
		Code:    types.NotAuthorized,
		Message: "Not Authorized",
	}

	if len(message) > 0 {
		vo.Message = message[0]
	}

	c.JSON(http.StatusUnauthorized, vo)
}