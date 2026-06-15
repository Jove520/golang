package handler

import (
	"context"
	"errors"
	"fmt"
	"geekai-rebuild/core"
	"geekai-rebuild/core/middleware"
	"geekai-rebuild/core/types"
	"geekai-rebuild/store/model"
	"geekai-rebuild/store/vo"
	"geekai-rebuild/utils"
	"geekai-rebuild/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserHandler struct {
	App          *core.AppServer
	DB           *gorm.DB
	Redis        *redis.Client
	SystemConfig *core.SystemConfigCache
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	CaptchaId   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaId   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

func NewUserHandler(app *core.AppServer, db *gorm.DB, redisClient *redis.Client, systemConfig *core.SystemConfigCache) *UserHandler {
	return &UserHandler{
		App:          app,
		DB:           db,
		Redis:        redisClient,
		SystemConfig: systemConfig,
	}
}

func (h *UserHandler) RegisterRoutes() {
	api := h.App.Engine.Group("/api/user")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	auth := api.Group("")
	auth.Use(middleware.UserAuthMiddleware(h.App.Config, h.Redis))
	auth.GET("/profile", h.Profile)
	auth.POST("/logout", h.Logout)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.CaptchaId = strings.TrimSpace(req.CaptchaId)
	req.CaptchaCode = strings.TrimSpace(req.CaptchaCode)

	if req.Username == "" || req.Password == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if err := h.verifyCaptcha(req.CaptchaId, req.CaptchaCode); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	var exists model.User
	err := h.DB.Where("username = ?", req.Username).First(&exists).Error
	if err == nil {
		resp.ERROR(c, "用户名已存在")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.ERROR(c, err.Error())
		return
	}

	salt, err := utils.GenSalt(12)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	systemConfig := h.SystemConfig.Get()

	user := model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: utils.GenPassword(req.Password, salt),
		Salt:     salt,
		Status:   true,
		Power:    systemConfig.Base.InitPower,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var userVO vo.User
	_ = utils.CopyObject(&userVO, user)

	resp.SUCCESS(c, userVO)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.CaptchaId = strings.TrimSpace(req.CaptchaId)
	req.CaptchaCode = strings.TrimSpace(req.CaptchaCode)

	if req.Username == "" || req.Password == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if err := h.verifyCaptcha(req.CaptchaId, req.CaptchaCode); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var user model.User
	err := h.DB.Where("username = ?", req.Username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.ERROR(c, "用户名或密码错误")
		return
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if !user.Status {
		resp.ERROR(c, "用户已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, user.Salt, user.Password) {
		resp.ERROR(c, "用户名或密码错误")
		return
	}

	expireHours := h.App.Config.JWT.ExpireHours
	if expireHours <= 0 {
		expireHours = 24
	}

	ttl := time.Duration(expireHours) * time.Hour
	token, err := utils.GenToken(user.Id, h.App.Config.JWT.Secret, ttl)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	ctx := context.Background()
	sessionKey := fmt.Sprintf("users/%d", user.Id)
	if err := h.Redis.Set(ctx, sessionKey, token, ttl).Err(); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	var userVO vo.User
	_ = utils.CopyObject(&userVO, user)

	resp.SUCCESS(c, gin.H{
		"token": token,
		"user":  userVO,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userIdValue, exists := c.Get(middleware.ContextUserId)
	if !exists {
		resp.NotAuth(c)
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		resp.NotAuth(c)
		return
	}

	var user model.User
	if err := h.DB.First(&user, userId).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var userVO vo.User
	_ = utils.CopyObject(&userVO, user)

	resp.SUCCESS(c, userVO)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userIdValue, exists := c.Get(middleware.ContextUserId)
	if !exists {
		resp.NotAuth(c)
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		resp.NotAuth(c)
		return
	}

	sessionKey := fmt.Sprintf("users/%d", userId)
	if err := h.Redis.Del(context.Background(), sessionKey).Err(); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, "ok")
}

func (h *UserHandler) verifyCaptcha(captchaId, captchaCode string) error {
	captchaId = strings.TrimSpace(captchaId)
	captchaCode = strings.TrimSpace(captchaCode)

	if captchaId == "" || captchaCode == "" {
		return errors.New("captcha required")
	}

	key := fmt.Sprintf("captcha/%s", captchaId)

	code, err := h.Redis.GetDel(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("captcha expired")
	}

	if err != nil {
		return err
	}

	if code != captchaCode {
		return errors.New("invalid captcha")
	}

	return nil
}
