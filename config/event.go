package config

type EventConfig struct {
	EnabledListeners []string `mapstructure:"enabledListeners"`
}
