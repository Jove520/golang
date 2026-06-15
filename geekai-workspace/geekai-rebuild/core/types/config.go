package types

import "fmt"

type AppConfig struct {
	Listen   string
	MysqlDsn string
	Redis    RedisConfig
	JWT      JWTConfig
}

const ConfigKeySystem = "system"

type BaseConfig struct {
	Title           string   `json:"title"`
	Slogan          string   `json:"slogan"`
	EnabledRegister bool     `json:"enabled_register"`
	RegisterWays    []string `json:"register_ways"`
	InitPower       int      `json:"init_power"`
}

type SystemConfig struct {
	Base BaseConfig `json:"base"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Title:           "GeekAI Rebuild",
		Slogan:          "从零复现 GeekAI",
		EnabledRegister: true,
		RegisterWays:    []string{"username"},
		InitPower:       100,
	}
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
