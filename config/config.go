package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

const FileName = "settings.yaml"

type Config struct {
	Port            string `json:"port" yaml:"port"`
	FDBHost         string `json:"fdb_host" yaml:"fdb_host"`
	TMDBProxy       bool   `json:"tmdb_proxy" yaml:"tmdb_proxy"`
	TMDBBearerToken string `json:"tmdb_bearer_token" yaml:"tmdb_bearer_token"`
	TGBotToken      string `json:"telegram_bot_token" yaml:"telegram_bot_token"`
	TGHost          string `json:"telegram_api_host" yaml:"telegram_api_host"`
	TSHost          string `json:"torrserver_host" yaml:"torrserver_host"`
}

var (
	mu          sync.RWMutex
	current     Config
	currentPath string
)

func Defaults() Config {
	return Config{
		Port:      "8094",
		FDBHost:   strings.TrimSpace(os.Getenv("FDBHOST")),
		TGHost:    "http://127.0.0.1:8081",
		TSHost:    "http://127.0.0.1:8090",
		TMDBProxy: false,
	}
}

func DefaultPath(baseDir string) string {
	return filepath.Join(baseDir, FileName)
}

func Load(path string) (Config, error) {
	cfg := Defaults()
	buf, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, err
		}
		normalize(&cfg)
		Set(path, cfg)
		return cfg, Save(path, cfg)
	}
	if len(strings.TrimSpace(string(buf))) > 0 {
		if err := yaml.Unmarshal(buf, &cfg); err != nil {
			return Config{}, err
		}
	}
	if cfg.FDBHost == "" {
		cfg.FDBHost = strings.TrimSpace(os.Getenv("FDBHOST"))
	}
	normalize(&cfg)
	Set(path, cfg)
	return cfg, nil
}

func Save(path string, cfg Config) error {
	normalize(&cfg)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, buf, 0o600); err != nil {
		return err
	}
	Set(path, cfg)
	return nil
}

func Set(path string, cfg Config) {
	normalize(&cfg)
	mu.Lock()
	defer mu.Unlock()
	currentPath = path
	current = cfg
}

func Current() Config {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

func Path() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentPath
}

func normalize(cfg *Config) {
	defaults := Defaults()
	cfg.Port = firstNonEmpty(cfg.Port, defaults.Port)
	cfg.FDBHost = strings.TrimSpace(cfg.FDBHost)
	cfg.TMDBBearerToken = strings.TrimSpace(cfg.TMDBBearerToken)
	cfg.TGBotToken = strings.TrimSpace(cfg.TGBotToken)
	cfg.TGHost = firstNonEmpty(cfg.TGHost, defaults.TGHost)
	cfg.TSHost = firstNonEmpty(cfg.TSHost, defaults.TSHost)
}

func firstNonEmpty(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	return fallback
}
