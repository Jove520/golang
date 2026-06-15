package model

import "time"

type ChatApp struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null;comment:App name" json:"name"`
	Context   string    `gorm:"column:context;type:text;comment:Preset context" json:"context"`
	Enabled   bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1;comment:Enabled flag" json:"enabled"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
