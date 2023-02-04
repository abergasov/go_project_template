package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AppPort        int       `yaml:"app_port"`
	MigratesFolder string    `yaml:"migrates_folder"`
	ConfigDB       DBConf    `yaml:"conf_db"`
	ConfigGraph    GraphConf `yaml:"conf_graph"`
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
