package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig `mapstructure:"server"`
	Dog          DogConfig    `mapstructure:"dog"`
	Cat          CatConfig    `mapstructure:"cat"`
	MaxOffenders int          `mapstructure:"maxOffenders"`
	CloseTimeout int          `mapstructure:"closeTimeout"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DogConfig struct {
	ConnString string `mapstructure:"connString"`
}

type CatConfig struct {
	ConnString string `mapstructure:"connString"`
}

func MustLoad() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config parsing error: %w", err)
	}

	if err := validate(cfg); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
