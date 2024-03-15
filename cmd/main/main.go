package main

import (
	//	"flag"
	"log"
	"os"
	"path/filepath"
	"torrsru/config"
	"torrsru/db/search"
	"torrsru/db/sync"
	"torrsru/web"
	"torrsru/web/global"
)

func main() {
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

	web.Start(port)
}
