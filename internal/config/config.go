package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Gate struct {
	Name     string `json:"name"`
	Mnemonic string `json:"mnemonic"`
}

type Config struct {
	Server ServerConfig `json:"config"`
	Gates  []Gate       `json:"gates"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// WALLET_MNEMONIC_<ИМЯ_ШЛЮЗА_UPPERCASE> перекрывает значение из JSON.
	for i := range cfg.Gates {
		envKey := "WALLET_MNEMONIC_" + strings.ToUpper(cfg.Gates[i].Name)
		if val := os.Getenv(envKey); val != "" {
			cfg.Gates[i].Mnemonic = val
		}
	}

	return &cfg, nil
}

func FindGate(cfg *Config, name string) *Gate {
	for i := range cfg.Gates {
		if cfg.Gates[i].Name == name {
			return &cfg.Gates[i]
		}
	}
	return nil
}
