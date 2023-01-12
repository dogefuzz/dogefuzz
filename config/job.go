package config

type JobConfig struct {
	EnabledJobs []string `mapstructure:"enabledJobs"`
}
