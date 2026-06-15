package model

type Config struct {
	Id    uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"column:name;type:varchar(50);uniqueIndex;not null;comment:配置名称" json:"name"`
	Value string `gorm:"column:value;type:text;not null;comment:配置 JSON" json:"value"`
}
