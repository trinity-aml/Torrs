package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"torrsru/config"
	"torrsru/models/fdb"
	"torrsru/web/global"
)

var (
	mu     sync.Mutex
	isSync bool
)

func StartSync() {
	for !global.Stopped {
		syncDB()
		time.Sleep(time.Minute * 20)
	}
}

func syncDB() {
	flag := ""
	mu.Lock()
	if isSync {
		mu.Unlock()
		return
	}
	isSync = true
	defer func() { isSync = false }()

	filetime := GetFileTime()
	lastft := filetime

	mu.Unlock()
	start := time.Now()
	gcCount := 0
	for {
		ftstr := strconv.FormatInt(filetime, 10)
		log.Println("Fetch:", ftstr)
		_, err := os.Stat("config.yml")
		if err != nil {
			if os.IsNotExist(err) {
				flag = config.Bypass + ftstr
			}
		} else {
			flag = config.ReadConfigParser("JacRed") + "/sync/fdb/torrents?time=" + ftstr
		}
		resp, err2 := http.Get(flag)
		if err2 != nil {
			log.Fatal("Error connect to fdb:", err)
			return
		}
		var js *fdb.FDBRequest
		err = json.NewDecoder(resp.Body).Decode(&js)
		if err != nil {
			log.Fatal("Error decode json:", err)
			return
		}
		resp.Body.Close()

		err = saveTorrent(js.Collections)
		if err != nil {
			log.Fatal("Error save torrents:", err)
			return
		}

		err = SetFileTime(filetime)
		if err != nil {
			log.Fatal("Error set ftime:", err)
			return
		}

		torrents := 0
		for _, col := range js.Collections {
			if col.Value.FileTime > filetime {
				filetime = col.Value.FileTime
			}
			torrents += len(col.Value.Torrents)
		}

		t := time.Unix(ft2sec(filetime), 0)
		log.Println("Save:", t.Format("2006-01-02 15:04:05"), ", Torrents:", torrents)

		if !js.Nextread {
			break
		}
		js = nil
		gcCount++
		if gcCount > 10 {
			runtime.GC()
			gcCount = 0
		}
	}
	if lastft != filetime {
		global.IsUpdateIndex = true
	}
	fmt.Println("End sync", time.Since(start))
}

func getHash(magnet string) string {
	pos := strings.Index(magnet, "btih:")
	if pos == -1 {
		return ""
	}
	magnet = magnet[pos+5:]
	pos = strings.Index(magnet, "&")
	if pos == -1 {
		return strings.ToLower(magnet)
	}
	return strings.ToLower(magnet[:pos])
}

func ft2sec(ft int64) int64 {
	//#define TICKS_PER_SECOND 10000000
	//#define EPOCH_DIFFERENCE 11644473600LL
	return ft/10000000 - 11644473600
}
