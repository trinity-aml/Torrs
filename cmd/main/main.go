package main

import (
	"fmt"
	"torrsru/version"
	//	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"torrsru/config"
	"torrsru/db/search"
	"torrsru/db/sync"
	"torrsru/web"
	"torrsru/web/global"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("=========== START ===========")
	fmt.Println("Torrs", version.Version+",", runtime.Version()+",", "CPU Num:", runtime.NumCPU())
	_, err := os.Stat("config.yml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Warning! Config file config.yml does not exist") // это_true
		}
	}

	pwd := filepath.Dir(os.Args[0])
	pwd, _ = filepath.Abs(pwd)
	log.Println("PWD:", pwd)
	global.PWD = pwd
	sync.Init()
	search.UpdateIndex()

	port := config.ReadConfigParser("Port")

	if port == "" {
		port = "8093"
	}

	fmt.Println("Port:", port)

	web.Start(port)
}
