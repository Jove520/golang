package service

import (
	"context"
	"errors"
	"time"

	"geekai-rebuild/store/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNoAvailableApiKey = errors.New("no available api key")

type ApiKeyService struct {
	DB *gorm.DB
}

func NewApiKeyService(db *gorm.DB) *ApiKeyService {
	return &ApiKeyService{
		DB: db,
	}
}

func (s *ApiKeyService) PickKey(ctx context.Context, keyType string) (*model.ApiKey, error) {
	var apiKey model.ApiKey
	now := time.Now()

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("type = ? AND enabled = ?", keyType, true).
			Order("last_used_at IS NULL DESC").
			Order("last_used_at ASC").
			Order("id ASC").
			First(&apiKey).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoAvailableApiKey
		}

		if err != nil {
			return err
		}

		if err := tx.Model(&model.ApiKey{}).
			Where("id = ?", apiKey.Id).
			Update("last_used_at", now).Error; err != nil {
			return err
		}

		apiKey.LastUsedAt = &now
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}
