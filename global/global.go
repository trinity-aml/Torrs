package global

import "github.com/gin-gonic/gin"

var (
	Route     *gin.Engine
	Stopped   = false
	PWD       = ""
	FDBHost   = ""
	TMDBProxy = false
	TMDBToken = ""
	TSHost    = ""

	SendFromWeb func(initData, msg string) error
)
