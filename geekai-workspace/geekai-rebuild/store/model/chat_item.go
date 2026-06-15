package model

import "time"

type ChatItem struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"column:user_id;type:int unsigned;not null;index;comment:User ID" json:"user_id"`
	Title     string    `gorm:"column:title;type:varchar(255);not null;default:'';comment:Chat title" json:"title"`
	Model     string    `gorm:"column:model;type:varchar(100);not null;default:'';comment:Model value" json:"model"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
