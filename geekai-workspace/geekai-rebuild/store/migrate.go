package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"geekai-rebuild/core/types"
	"geekai-rebuild/store/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AutoMigrateAndSeed(db *gorm.DB, log *zap.SugaredLogger) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Config{},
		&model.ChatModel{},
		&model.ApiKey{},
		&model.ChatApp{},
		&model.ChatItem{},
		&model.ChatMessage{},
	); err != nil {
		return err
	}

	if err := ensureChatIntegerColumnTypes(db); err != nil {
		return err
	}

	if err := seedSystemConfig(db); err != nil {
		return err
	}

	if err := seedChatModels(db); err != nil {
		return err
	}

	if err := seedApiKeys(db); err != nil {
		return err
	}
	log.Info("database migration completed")
	return nil
}

func seedSystemConfig(db *gorm.DB) error {
	var config model.Config
	err := db.Where("name = ?", types.ConfigKeySystem).First(&config).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	value, err := json.Marshal(types.DefaultBaseConfig())
	if err != nil {
		return err
	}

	return db.Create(&model.Config{
		Name:  types.ConfigKeySystem,
		Value: string(value),
	}).Error
}

func ensureChatIntegerColumnTypes(db *gorm.DB) error {
	columns := []struct {
		table      string
		column     string
		definition string
	}{
		{
			table:      db.Config.NamingStrategy.TableName("chat_models"),
			column:     "power",
			definition: "int NOT NULL DEFAULT 1 COMMENT 'Power cost'",
		},
		{
			table:      db.Config.NamingStrategy.TableName("chat_models"),
			column:     "max_tokens",
			definition: "int NOT NULL DEFAULT 4096 COMMENT 'Max output tokens'",
		},
		{
			table:      db.Config.NamingStrategy.TableName("chat_models"),
			column:     "max_context",
			definition: "int NOT NULL DEFAULT 8192 COMMENT 'Max context tokens'",
		},
		{
			table:      db.Config.NamingStrategy.TableName("chat_messages"),
			column:     "tokens",
			definition: "int NOT NULL DEFAULT 0 COMMENT 'Token count'",
		},
	}

	for _, column := range columns {
		if err := ensureColumnType(db, column.table, column.column, "int", column.definition); err != nil {
			return err
		}
	}

	return nil
}

func ensureColumnType(db *gorm.DB, table, column, expectedType, definition string) error {
	var columnType string
	if err := db.Raw(
		"SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?",
		table,
		column,
	).Scan(&columnType).Error; err != nil {
		return err
	}

	if columnType == expectedType {
		return nil
	}

	sql := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s", table, column, definition)
	return db.Exec(sql).Error
}

func seedChatModels(db *gorm.DB) error {
	var count int64

	if err := db.Model(&model.ChatModel{}).
		Where("value = ?", "gpt-4o-mini").
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return db.Create(&model.ChatModel{
		Name:       "GPT-4o Mini",
		Value:      "gpt-4o-mini",
		Power:      1,
		MaxTokens:  4096,
		MaxContext: 8192,
		Enabled:    true,
	}).Error
}

func seedApiKeys(db *gorm.DB) error {
	var count int64

	if err := db.Model(&model.ApiKey{}).
		Where("type = ? AND value = ?", "openai", "CHANGE_ME").
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	now := time.Now()
	return db.Model(&model.ApiKey{}).Create(map[string]any{
		"type":       "openai",
		"value":      "CHANGE_ME",
		"api_url":    "https://api.openai.com/v1/chat/completions",
		"enabled":    false,
		"created_at": now,
		"updated_at": now,
	}).Error
}
