package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"torrsru/config"
	"torrsru/db"
	"torrsru/global"
	"torrsru/tgbot"
	"torrsru/utils"
	"torrsru/version"
	"torrsru/web"
)

func main() {
	port := ""
	ri := false
	jac := ""

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("=========== START ===========")
	fmt.Println("Torrs", version.Version+",", runtime.Version()+",", "CPU Num:", runtime.NumCPU())

	var args struct {
		Port         string `default:"8094" arg:"-p" help:"port for http"`
		RebuildIndex bool   `default:"false" arg:"-r" help:"rebuild index and exit"`
		JacRed       string `default:"" arg:"-j" help:"alternative address JacRed"`
		TMDBProxy    bool   `default:"false" arg:"--tmdb" help:"proxy for TMDB"`
		TGBotToken   string `default:"" arg:"--token" help:"telegram bot token"`
		TGHost       string `default:"http://127.0.0.1:8081" arg:"--tgapi" help:"telegram api host"`
		TSHost       string `default:"http://127.0.0.1:8090" arg:"--ts" help:"TorrServer host"`
	}
	arg.MustParse(&args)

	if config.ReadConfigParser("Port") != "" {
		port = config.ReadConfigParser("Port")
	} else {
		port = args.Port
	}
	if config.ReadConfigParser("Rebuild") == "true" {
		ri = true
	} else {
		ri = args.RebuildIndex
	}
	if config.ReadConfigParser("JacRed") != "" {
		jac = config.ReadConfigParser("JacRed")
	} else {
		jac = args.JacRed
	}

	fmt.Println("Port:", port)
	fmt.Println("Rebuild index of base:", ri)
	fmt.Println("Alternative JacRed address:", jac)

	if jac == "" {
		jac = "http://62.112.8.193:9117"
	}
	os.Setenv("JacRed", jac)

	pwd := filepath.Dir(os.Args[0])
	pwd, _ = filepath.Abs(pwd)
	log.Println("PWD:", pwd)
	global.PWD = pwd

	global.TMDBProxy = args.TMDBProxy
	global.TSHost = args.TSHost

	db.Init()

	if ri {
		err := db.RebuildIndex()
		if err != nil {
			log.Println("Rebuild index error:", err)
		} else {
			log.Println("Rebuild index success")
		}
		return
	}

	if args.TGBotToken != "" {
		if args.TGHost == "" {
			log.Println("Error telegram host is empty. Telegram api bot need for upload 2gb files")
			os.Exit(1)
		}
		err := tgbot.Start(args.TGBotToken, args.TGHost)
		if err != nil {
			log.Println("Start Telegram bot error:", err)
			os.Exit(1)
		}
	}

	for !utils.TestServer(jac) {
		time.Sleep(10 * time.Second)
	}

	web.Start(port)
}
