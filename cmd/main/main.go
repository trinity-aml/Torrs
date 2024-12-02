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
	"torrsru/utils"
	"torrsru/version"
	"torrsru/web"
	"torrsru/web/global"
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

	os.Setenv("JacRed", jac)

	pwd := filepath.Dir(os.Args[0])
	pwd, _ = filepath.Abs(pwd)
	log.Println("PWD:", pwd)
	global.PWD = pwd

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

	if jac == "" {
		jac = "http://62.112.8.193:9117"
	}

	for !utils.TestServer(jac) {
		time.Sleep(10 * time.Second)
	}

	web.Start(port)
}
