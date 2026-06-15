package model

import "time"

type ChatModel struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null;comment:Model name" json:"name"`
	Value      string    `gorm:"column:value;type:varchar(100);uniqueIndex;not null;comment:Model value" json:"value"`
	Power      int       `gorm:"column:power;type:int;not null;default:1;comment:Power cost" json:"power"`
	MaxTokens  int       `gorm:"column:max_tokens;type:int;not null;default:4096;comment:Max output tokens" json:"max_tokens"`
	MaxContext int       `gorm:"column:max_context;type:int;not null;default:8192;comment:Max context tokens" json:"max_context"`
	Enabled    bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1;comment:Enabled flag" json:"enabled"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
