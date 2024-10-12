package httpx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeyCheek888/go-std/httpx/websocket"
	websockerManager "github.com/joeyCheek888/go-std/httpx/websocket/manager"
	"github.com/joeyCheek888/go-std/log"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	conf *Config
	http.Server
	Router           *gin.Engine
	websocketManager *websockerManager.Manager
}

func NewServer(conf *Config) *Server {

	conf.Check()

	s := &Server{
		conf:             conf,
		Server:           http.Server{},
		Router:           newRouter(),
		websocketManager: nil,
	}

	return s
}

// swag init -g ./controller.go -d ./internal/controller/customer_support/business --parseDependency -o ./docs/customer_support/business

// EnableWebsocket 开启websocket
func (s *Server) EnableWebsocket(path string) {
	log.Logger.Info("启用websocket", zap.String("path", path))

	s.websocketManager = websockerManager.NewManager()
	go s.websocketManager.Start()

	websocket.Upgrade(path, s.websocketManager, s.Router)
}

func (s *Server) WebsocketManger() *websockerManager.Manager {

	if s.websocketManager == nil {
		log.Logger.Error("websocket not initialized")
		panic("websocket not initialized!")
	}

	return s.websocketManager
}

func (s *Server) Start() {

	address := fmt.Sprintf("%s:%s", s.conf.Addr, s.conf.Port)

	log.Logger.Info("启动HTTP服务", log.String("地址", address))

	s.Server.Addr = address
	s.Server.Handler = s.Router

	err := s.Server.ListenAndServe()
	if err != nil {
		log.Logger.Error("启动HTTP服务失败", log.Error(err))
		return
	}

}

func (s *Server) Stop() {
	err := s.Server.Shutdown(nil)
	if err != nil {
		return
	}
}
