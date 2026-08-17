package config

import (
	"fmt"
	"go_project_template/internal/utils/config"
	"path/filepath"
)

type AppConfig struct {
	AppPort         int       `yaml:"app_port" env:"APP_PORT" flag:"app_port" default:"8000"`
	EnableTelemetry bool      `yaml:"enable_telemetry" env:"ENABLE_TELEMETRY" flag:"enable_telemetry" default:"true"`
	MigratesFolder  string    `yaml:"migrates_folder" env:"MIGRATES_FOLDER" flag:"migrates_folder" default:"migrations"`
	ConfigDB        DBConf    `yaml:"conf_db"`
	ConfigGraph     GraphConf `yaml:"conf_graph"`
}

type GraphConf struct {
	Address        string `yaml:"address" json:"address,omitempty"`
	Port           string `yaml:"port" json:"port,omitempty"`
	User           string `yaml:"user" json:"user,omitempty"`
	Pass           string `yaml:"pass" json:"pass,omitempty"`
	DBName         string `yaml:"db_name" json:"db_name,omitempty"`
	MaxConnections int    `yaml:"max_connections" json:"max_connections,omitempty"`
}

type DBConf struct {
	Address        string `yaml:"address" env:"DB_ADDRESS" flag:"db_address" default:"localhost"`
	Port           string `yaml:"port" env:"DB_PORT" flag:"db_port" default:"5432"`
	User           string `yaml:"user" env:"DB_USER" flag:"db_user" default:"user"`
	Pass           string `yaml:"pass" env:"DB_PASS" flag:"db_pass"`
	DBName         string `yaml:"db_name" env:"DB_NAME" flag:"db_name" default:"dbname"`
	MaxConnections int    `yaml:"max_connections" env:"DB_MAX_CONNECTIONS" flag:"db_max_connections" default:"10"`
}

func InitConf(confFile *string) (*AppConfig, error) {
	options := make([]config.Option, 0, 1)
	if confFile != nil && *confFile != "" {
		options = append(options, config.WithYAML(filepath.Clean(*confFile)))
	}

	var cfg AppConfig
	if err := config.Load(&cfg, options...); err != nil {
		return nil, fmt.Errorf("error load config: %w", err)
	}
	return &cfg, nil
}
