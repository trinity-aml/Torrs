package config

import "github.com/ilyakaznacheev/cleanenv"

const Bypass = "http://62.112.8.193:9117/sync/fdb/torrents?time="

type ConfParser struct {
	Port   string `yaml:"port"`
	JacRed string `yaml:"jacred"`
}

var cfg ConfParser

func ReadConfigParser(vars string) string {
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err == nil {
		switch {
		case vars == "Port":
			return cfg.Port
		case vars == "JacRed":
			return cfg.JacRed
		}
	}
	return ""
}
