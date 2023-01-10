package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	StorageFolder  string         `mapstructure:"storageFolder"`
	DatabaseName   string         `mapstructure:"databaseName"`
	ServerPort     int            `mapstructure:"serverPort"`
	GethConfig     GethConfig     `mapstructure:"geth"`
	FuzzerConfig   FuzzerConfig   `mapstructure:"fuzzer"`
	VandalConfig   VandalConfig   `mapstructure:"vandal"`
	JobConfig      JobConfig      `mapstructure:"job"`
	EventConfig    EventConfig    `mapstructure:"event"`
	ReporterConfig ReporterConfig `mapstructure:"reporter"`
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
