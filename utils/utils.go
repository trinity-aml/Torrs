package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func TestServer(u string) bool {
	if u == "" {
		return false
	}
	client := http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
		return false
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		if resp.StatusCode == 200 {
			fmt.Println(u + ": " + "ok")
		} else {
			fmt.Println("Server link ", u, " returned error code: ", resp.StatusCode)
		}
	}
	defer resp.Body.Close()
	return true
}
