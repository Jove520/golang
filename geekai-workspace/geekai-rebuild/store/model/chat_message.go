package model

import "time"

type ChatMessage struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     uint      `gorm:"column:user_id;type:int unsigned;not null;index;comment:User ID" json:"user_id"`
	ChatItemId uint      `gorm:"column:chat_item_id;type:int unsigned;not null;index;comment:Chat item ID" json:"chat_item_id"`
	Role       string    `gorm:"column:role;type:varchar(20);not null;comment:Message role" json:"role"`
	Content    string    `gorm:"column:content;type:text;comment:Message content" json:"content"`
	Tokens     int       `gorm:"column:tokens;type:int;not null;default:0;comment:Token count" json:"tokens"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
