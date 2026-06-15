package model

import "time"

type ApiKey struct {
	Id         uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Type       string     `gorm:"column:type;type:varchar(30);not null;index;comment:Key type" json:"type"`
	Value      string     `gorm:"column:value;type:varchar(255);not null;comment:API Key" json:"value"`
	ApiUrl     string     `gorm:"column:api_url;type:varchar(255);not null;comment:API URL" json:"api_url"`
	LastUsedAt *time.Time `gorm:"column:last_used_at;type:datetime;comment:Last used time" json:"last_used_at"`
	Enabled    bool       `gorm:"column:enabled;type:tinyint(1);not null;default:1;comment:Enabled flag" json:"enabled"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
