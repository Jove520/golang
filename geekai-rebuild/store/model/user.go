package model

import "time"

type User struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(30);uniqueIndex;not null;comment:用户名" json:"username"`
	Nickname  string    `gorm:"column:nickname;type:varchar(30);not null;comment:昵称" json:"nickname"`
	Password  string    `gorm:"column:password;type:char(64);not null;comment:密码" json:"password"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像" json:"avatar"`
	Salt      string    `gorm:"column:salt;type:char(12);not null;comment:密码盐" json:"salt"`
	Power     int       `gorm:"column:power;type:int;not null;default:0;comment:剩余算力" json:"power"`
	Status    bool      `gorm:"column:status;type:tinyint(1);not null;default:1;comment:当前状态" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}
