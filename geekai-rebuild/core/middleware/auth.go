package middleware

import (
	"context"
	"fmt"
	"strings"

	"geekai-rebuild/core/types"
	"geekai-rebuild/utils"
	"geekai-rebuild/utils/resp"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const ContextUserId = "user_id"

func UserAuthMiddleware(appConfig *types.AppConfig, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			resp.NotAuth(c, "missing token")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString, appConfig.JWT.Secret)
		if err != nil {
			resp.NotAuth(c, "invalid token")
			c.Abort()
			return
		}

		sessionKey := fmt.Sprintf("users/%d", claims.UserId)
		exists, err := redisClient.Exists(context.Background(), sessionKey).Result()
		if err != nil {
			resp.NotAuth(c, err.Error())
			c.Abort()
			return
		}
		if exists == 0 {
			resp.NotAuth(c, "session expired")
			c.Abort()
			return
		}

		c.Set(ContextUserId, claims.UserId)
		c.Next()
	}
}

func getBearerToken(authorization string) string {
	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return ""
	}

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
