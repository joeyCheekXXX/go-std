package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/joeyCheek888/go-std/httpx/websocket/handler"
	"testing"
)

func TestWebsocket(t *testing.T) {

	conf := &Config{
		Addr:             "",
		Port:             "",
		RouterPrefix:     "",
		Mode:             "",
		EnableSwaggerDoc: false,
	}

	s := NewServer(conf)
	defer s.Stop()

	s.EnableWebsocket("/test/ws", func(context *gin.Context) {
		handler.Upgrade(context.Writer, context.Request)
	})

}
