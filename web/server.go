package web

import (
	"log"
	"strings"
	ss "sync"
	"time"
	"torrsru/db"
	"torrsru/global"
	"torrsru/web/api"
	"torrsru/web/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(port string) {
	go db.StartSync()

	//gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "Accept", "Authorization"}

	global.Route = gin.New()
	global.Route.Use(gin.Recovery(), cors.New(corsCfg), blockUsers())
	static.RouteEmbedFiles(global.Route)
	api.SetRoutes(global.Route)

	err := global.Route.Run(":" + port)
	if err != nil {
		log.Println("Error start server:", err)
	}

	global.Stopped = true
}

func blockUsers() gin.HandlerFunc {
	var mu ss.Mutex
	return func(c *gin.Context) {
		referer := strings.ToLower(c.Request.Referer())
		useragent := strings.ToLower(c.Request.UserAgent())

		if strings.Contains(referer, "lamp") || strings.Contains(useragent, "lamp") {
			mu.Lock()
			c.Next()
			time.Sleep(time.Millisecond * 300)
			mu.Unlock()
			return
		}

		c.Next()
	}
}
