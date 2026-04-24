package settings

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"torrsru/config"
	"torrsru/global"
)

const secretMask = "********"

type settingsPayload struct {
	Port               string `json:"port"`
	FDBHost            string `json:"fdb_host"`
	TMDBProxy          bool   `json:"tmdb_proxy"`
	TMDBBearerToken    string `json:"tmdb_bearer_token"`
	TMDBBearerTokenSet bool   `json:"tmdb_bearer_token_set"`
	TGBotToken         string `json:"telegram_bot_token"`
	TGBotTokenSet      bool   `json:"telegram_bot_token_set"`
	TGHost             string `json:"telegram_api_host"`
	TSHost             string `json:"torrserver_host"`
}

type settingsResponse struct {
	Path     string          `json:"path"`
	Settings settingsPayload `json:"settings"`
	Message  string          `json:"message,omitempty"`
}

func Page(c *gin.Context) {
	if !localOnly(c) {
		return
	}
	c.HTML(http.StatusOK, "settings.go.html", gin.Mode())
}

func Get(c *gin.Context) {
	if !localOnly(c) {
		return
	}
	c.JSON(http.StatusOK, response(config.Current(), ""))
}

func Save(c *gin.Context) {
	if !localOnly(c) {
		return
	}
	var req settingsPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cfg := config.Current()
	cfg.Port = req.Port
	cfg.FDBHost = req.FDBHost
	cfg.TMDBProxy = req.TMDBProxy
	cfg.TGHost = req.TGHost
	cfg.TSHost = req.TSHost
	if req.TMDBBearerToken != secretMask {
		cfg.TMDBBearerToken = req.TMDBBearerToken
	}
	if req.TGBotToken != secretMask {
		cfg.TGBotToken = req.TGBotToken
	}

	path := config.Path()
	if path == "" {
		path = config.DefaultPath(global.PWD)
	}
	if err := config.Save(path, cfg); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cfg = config.Current()
	applyRuntime(cfg)
	c.JSON(http.StatusOK, response(cfg, "Настройки сохранены"))
}

func response(cfg config.Config, msg string) settingsResponse {
	return settingsResponse{
		Path: config.Path(),
		Settings: settingsPayload{
			Port:               cfg.Port,
			FDBHost:            cfg.FDBHost,
			TMDBProxy:          cfg.TMDBProxy,
			TMDBBearerToken:    maskSecret(cfg.TMDBBearerToken),
			TMDBBearerTokenSet: cfg.TMDBBearerToken != "",
			TGBotToken:         maskSecret(cfg.TGBotToken),
			TGBotTokenSet:      cfg.TGBotToken != "",
			TGHost:             cfg.TGHost,
			TSHost:             cfg.TSHost,
		},
		Message: msg,
	}
}

func maskSecret(value string) string {
	if value == "" {
		return ""
	}
	return secretMask
}

func applyRuntime(cfg config.Config) {
	global.FDBHost = cfg.FDBHost
	global.TMDBProxy = cfg.TMDBProxy
	global.TMDBToken = cfg.TMDBBearerToken
	global.TSHost = cfg.TSHost
}

func localOnly(c *gin.Context) bool {
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		host = c.Request.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		c.AbortWithStatus(http.StatusForbidden)
		return false
	}
	return true
}
