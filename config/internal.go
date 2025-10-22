package config

import "fmt"

func validate(cfg Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be 1-65535)", cfg.Server.Port)
	}
	if cfg.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if cfg.Dog.ConnString == "" {
		return fmt.Errorf("dog connString cannot be empty")
	}
	if cfg.Cat.ConnString == "" {
		return fmt.Errorf("cat connString cannot be empty")
	}
	if cfg.MaxOffenders <= 0 {
		return fmt.Errorf("maxOffenders must be > 0")
	}
	if cfg.CloseTimeout < 0 {
		return fmt.Errorf("closeTimeout cannot be negative")
	}
	return nil
}
