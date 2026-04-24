package tmdb

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"torrsru/global"
)

func TMDBAPI(c *gin.Context) {
	token, ok := tmdbToken(c)
	if !ok {
		return
	}
	url := "https://api.themoviedb.org/" + strings.TrimPrefix(c.Request.RequestURI, "/tmdb/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", bearerHeader(token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	defer res.Body.Close()

	c.Header("Cache-Control", "public, max-age=604800") //1 week
	c.Header("Etag", res.Header.Get("Etag"))
	c.DataFromReader(res.StatusCode, res.ContentLength, "application/javascript; charset=utf-8", res.Body, nil)
}

func TMDBIMG(c *gin.Context) {
	token, ok := tmdbToken(c)
	if !ok {
		return
	}
	//url := "https://imagetmdb.com/" + strings.TrimPrefix(c.Request.RequestURI, "/tmdbimg/")
	url := "https://image.tmdb.org/" + strings.TrimPrefix(c.Request.RequestURI, "/tmdbimg/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	req.Header.Add("Authorization", bearerHeader(token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	defer res.Body.Close()

	c.Header("Cache-Control", "public, max-age=604800") //1 week
	c.Header("Etag", res.Header.Get("Etag"))
	c.DataFromReader(res.StatusCode, res.ContentLength, res.Header.Get("Content-Type"), res.Body, nil)
}

func tmdbToken(c *gin.Context) (string, bool) {
	if !global.TMDBProxy {
		c.Status(http.StatusNotFound)
		return "", false
	}
	token := strings.TrimSpace(global.TMDBToken)
	if token == "" {
		c.String(http.StatusServiceUnavailable, "TMDB bearer token is not configured")
		return "", false
	}
	return token, true
}

func bearerHeader(token string) string {
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token
	}
	return "Bearer " + token
}
