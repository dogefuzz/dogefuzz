package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	StorageFolder string `mapstructure:"STORAGE_FOLDER"`
	DatabaseName  string `mapstructure:"DATABASE_NAME"`
	ServerPort    int    `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
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
