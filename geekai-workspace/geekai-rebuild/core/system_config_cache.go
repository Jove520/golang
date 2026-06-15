package core

import (
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"

	"geekai-rebuild/core/types"
	"geekai-rebuild/store/model"

	"gorm.io/gorm"
)

type SystemConfigCache struct {
	db    *gorm.DB
	value atomic.Value // stores types.SystemConfig
}

func NewSystemConfigCache(db *gorm.DB) (*SystemConfigCache, error) {
	cache := &SystemConfigCache{db: db}
	if err := cache.Reload(); err != nil {
		return nil, err
	}
	return cache, nil
}

func (c *SystemConfigCache) Reload() error {
	var row model.Config
	if err := c.db.Where("name = ?", types.ConfigKeySystem).First(&row).Error; err != nil {
		return err
	}

	cfg, err := decodeSystemConfig(row.Value)
	if err != nil {
		return err
	}

	c.value.Store(cfg)
	return nil
}

func (c *SystemConfigCache) Get() types.SystemConfig {
	if value := c.value.Load(); value != nil {
		return value.(types.SystemConfig)
	}

	return types.SystemConfig{Base: types.DefaultBaseConfig()}
}

func decodeSystemConfig(raw string) (types.SystemConfig, error) {
	if strings.TrimSpace(raw) == "" {
		return types.SystemConfig{Base: types.DefaultBaseConfig()}, nil
	}

	var cfg types.SystemConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err == nil && cfg.Base.Title != "" {
		return cfg, nil
	}

	var base types.BaseConfig
	if err := json.Unmarshal([]byte(raw), &base); err != nil {
		return types.SystemConfig{}, err
	}
	if base.Title == "" {
		return types.SystemConfig{}, errors.New("invalid system config")
	}

	return types.SystemConfig{Base: base}, nil
}
