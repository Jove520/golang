package core

import (
	"bytes"
	"os"

	"geekai-rebuild/core/types"

	"github.com/BurntSushi/toml"
)

func NewDefaultConfig() *types.AppConfig {
	return &types.AppConfig{
		Listen:   "0.0.0.0:5678",
		MysqlDsn: "root:root@tcp(127.0.0.1:3306)/geekai?charset=utf8mb4&parseTime=True&loc=Local",
		Redis: types.RedisConfig{
			Host:     "127.0.0.1",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		JWT: types.JWTConfig{
			Secret:      "geekai-rebuild-secret",
			ExpireHours: 24,
		},
	}
}

func LoadConfig(configFile string) (*types.AppConfig, error) {
	config := NewDefaultConfig()

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			return config, SaveConfig(configFile, config)
		}

		return nil, err
	}

	if _, err := toml.DecodeFile(configFile, config); err != nil {
		return nil, err
	}

	return config, nil
}

func SaveConfig(configFile string, config *types.AppConfig) error {
	buf := new(bytes.Buffer)

	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	}

	return os.WriteFile(configFile, buf.Bytes(), 0644)
}
