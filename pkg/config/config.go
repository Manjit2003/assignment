package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var Config *APIConfig

type webCfg struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type dbCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Keyspace string `yaml:"keyspace"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type authCfg struct {
	JWTSecret       string `yaml:"jwt_secret"`
	TokenExpiration int    `yaml:"token_expiration"` // In minutes
}

type logCfg struct {
	Level    string `yaml:"level"`    // e.g., debug, info, warn, error
	Filename string `yaml:"filename"` // Log file path
}

type APIConfig struct {
	Web      webCfg  `yaml:"web"`
	Database dbCfg   `yaml:"database"`
	Auth     authCfg `yaml:"auth"`
	Logging  logCfg  `yaml:"logging"`
}

func LoadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config data: %v", err)
	}
	return nil
}
