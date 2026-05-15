package torr

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"torrsru/global"
	"torrsru/tgbot/torr/state"
)

type TorrentDetails struct {
	Title   string
	Size    string
	Date    time.Time
	Link    string
	Tracker string
	Peer    int
	Seed    int
	Magnet  string
}

var infoHTTPClient = &http.Client{Timeout: 30 * time.Second}

func GetTorrentInfo(hash string) (*state.TorrentStatus, error) {
	link := global.TSHost + "/stream?stat&link=" + url.QueryEscape(hash)
	resp, err := infoHTTPClient.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ti *state.TorrentStatus
	if err := json.NewDecoder(resp.Body).Decode(&ti); err != nil {
		return nil, err
	}
	return ti, nil
}
