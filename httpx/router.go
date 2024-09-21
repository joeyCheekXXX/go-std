package httpx

import "github.com/gin-gonic/gin"

func newRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	Router := gin.New()

	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	return Router
}
