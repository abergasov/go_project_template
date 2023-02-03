package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AppPort        int    `yaml:"app_port"`
	MigratesFolder string `yaml:"migrates_folder"`
	ConfigDB       DBConf `yaml:"conf_db"`
}

type DBConf struct {
	Address        string `yaml:"address"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	DBName         string `yaml:"db_name"`
	MaxConnections int    `yaml:"max_connections"`
}

func InitConf(confFile string) (*AppConfig, error) {
	file, err := os.Open(filepath.Clean(confFile))
	if err != nil {
		return nil, fmt.Errorf("error open config file: %w", err)
	}
	defer func() {
		if e := file.Close(); e != nil {
			log.Fatal("Error close config file", e)
		}
	}()

	var cfg AppConfig
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decode config file: %w", err)
	}

	return &cfg, nil
}
