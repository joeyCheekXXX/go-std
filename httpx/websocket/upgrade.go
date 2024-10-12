package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joeyCheek888/go-std/httpx/websocket/manager"
	"github.com/joeyCheek888/go-std/log"
	"net/http"
	"strings"
)

// Upgrade
//
//	@Description: 升级协议
//	@param path
//	@param _manager
//	@param router
func Upgrade(path string, _manager *manager.Manager, router *gin.Engine) {

	if strings.Contains(path, ":token") == false {
		log.Logger.Error(`websocket path 不包含 ":token"`)
		panic(`websocket path 不包含 ":token"`)
	}

	router.GET(path, func(context *gin.Context) {

		// 获取token
		token := context.Param("token")

		// 检查token
		ok := _manager.Event.CheckTokenFunc(token)
		if !ok {
			http.NotFound(context.Writer, context.Request)
			return
		}

		// 创建连接
		conn, err := (&websocket.Upgrader{
			HandshakeTimeout:  0,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			WriteBufferPool:   nil,
			Subprotocols:      nil,
			Error:             nil,
			CheckOrigin:       _manager.Event.CheckOriginFunc,
			EnableCompression: false,
		}).Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			http.NotFound(context.Writer, context.Request)
			return
		}

		// 注册连接
		_manager.Event.RegisterChan <- manager.NewMember(conn, "token")
	})

}
