package main

import (
	"github.com/alexflint/go-arg"
	"log"
	"os"
	"path/filepath"
	"torrsru/config"
	"torrsru/db"
	"torrsru/global"
	"torrsru/tgbot"
	"torrsru/web"
)

type argsConfig struct {
	Config          string  `arg:"-c,--config" help:"path to YAML settings file"`
	Port            *string `arg:"-p" help:"port for http"`
	RebuildIndex    bool    `arg:"-r" help:"rebuild index and exit"`
	TMDBProxy       *bool   `arg:"--tmdb" help:"proxy for TMDB"`
	TGBotToken      *string `arg:"--token" help:"telegram bot token"`
	TGHost          *string `arg:"--tgapi" help:"telegram api host"`
	TSHost          *string `arg:"--ts" help:"TorrServer host"`
	FDBHost         *string `arg:"--fdbhost" help:"FDB sync host"`
	TMDBBearerToken *string `arg:"--tmdb-token" help:"TMDB bearer token"`
}

func main() {
	pwd := filepath.Dir(os.Args[0])
	pwd, _ = filepath.Abs(pwd)

	args := argsConfig{Config: config.DefaultPath(pwd)}
	arg.MustParse(&args)
	args.Config = absPath(args.Config)

	log.Println("PWD:", pwd)
	global.PWD = pwd

	cfg, err := config.Load(args.Config)
	if err != nil {
		log.Fatalln("Error load config:", err)
	}

	cfg = applyArgs(cfg, args)
	config.Set(args.Config, cfg)

	global.FDBHost = cfg.FDBHost
	global.TMDBProxy = cfg.TMDBProxy
	global.TMDBToken = cfg.TMDBBearerToken
	global.TSHost = cfg.TSHost

	db.Init()

	if args.RebuildIndex {
		err := db.RebuildIndex()
		if err != nil {
			log.Println("Rebuild index error:", err)
		} else {
			log.Println("Rebuild index success")
		}
		return
	}

	if cfg.TGBotToken != "" {
		if cfg.TGHost == "" {
			log.Println("Error telegram host is empty. Telegram api bot need for upload 2gb files")
			os.Exit(1)
		}
		err := tgbot.Start(cfg.TGBotToken, cfg.TGHost)
		if err != nil {
			log.Println("Start Telegram bot error:", err)
			os.Exit(1)
		}
	}
	web.Start(cfg.Port)
}

func absPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func applyArgs(cfg config.Config, args argsConfig) config.Config {
	if args.Port != nil {
		cfg.Port = *args.Port
	}
	if args.TMDBProxy != nil {
		cfg.TMDBProxy = *args.TMDBProxy
	}
	if args.TGBotToken != nil {
		cfg.TGBotToken = *args.TGBotToken
	}
	if args.TGHost != nil {
		cfg.TGHost = *args.TGHost
	}
	if args.TSHost != nil {
		cfg.TSHost = *args.TSHost
	}
	if args.FDBHost != nil {
		cfg.FDBHost = *args.FDBHost
	}
	if args.TMDBBearerToken != nil {
		cfg.TMDBBearerToken = *args.TMDBBearerToken
	}
	return cfg
}
