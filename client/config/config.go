package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AnthropicConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type Config struct {
	Anthropic AnthropicConfig `mapstructure:"anthropic"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile("../cmd/conf/config.yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
