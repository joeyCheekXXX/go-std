package httpx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeyCheekXXX/go-std/httpx/websocket"
	"github.com/joeyCheekXXX/go-std/log"
	"net/http"
)

var _config *Config

type Server struct {
	http.Server
	Router          *gin.Engine
	enableWebsocket bool
}

func NewServer(conf *Config) *Server {

	_config = conf

	s := &Server{
		Server: http.Server{},
		Router: newRouter(),
	}
	gin.SetMode(_config.Mode)

	return s
}

// EnableWs 开启websocket
func (s *Server) EnableWs(path string, handler func(context *gin.Context)) {
	s.enableWebsocket = true
	s.Router.GET(path, handler)
}

func (s *Server) Start() {

	if _config.Addr == "" {
		_config.Addr = "0.0.0.0"
	}

	address := fmt.Sprintf("%s:%s", _config.Addr, _config.Port)

	log.Logger.Info("启动HTTP服务", log.String("地址", address))

	s.Server.Addr = address
	s.Server.Handler = s.Router

	if s.enableWebsocket {
		go websocket.Manager.Start()
	}

	err := s.Server.ListenAndServe()
	if err != nil {
		return
	}

}

func (s *Server) Stop() {
	err := s.Server.Shutdown(nil)
	if err != nil {
		return
	}
}
