package httpx

import "github.com/gin-gonic/gin"

func newRouter() *gin.Engine {

	Router := gin.New()

	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	return Router
}
