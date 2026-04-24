package api

import (
	"github.com/gin-gonic/gin"
	"torrsru/web/api/pages"
	"torrsru/web/api/settings"
	"torrsru/web/api/tgbot"
	"torrsru/web/api/tmdb"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", pages.MainPage)
	//r.GET("/robots.txt", pages.RobotsPage)
	r.GET("/search", pages.Search)
	r.GET("/settings", settings.Page)
	r.GET("/api/settings", settings.Get)
	r.POST("/api/settings", settings.Save)
	r.GET("/tmdb/*path", tmdb.TMDBAPI)
	r.GET("/tmdbimg/*path", tmdb.TMDBIMG)
	r.POST("/sendbot", tgbot.SendBot)
}
