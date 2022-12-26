package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	StorageFolder       string       `json:"storageFolder"`
	DatabaseName        string       `json:"databaseName"`
	ServerPort          int          `json:"serverPort"`
	GethConfig          GethConfig   `json:"geth"`
	FuzzerConfig        FuzzerConfig `json:"fuzzer"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.StorageFolder == "" {
		cfg.StorageFolder = path.Join(os.TempDir(), "dogefuzz")
	}

	return &cfg, nil
}
