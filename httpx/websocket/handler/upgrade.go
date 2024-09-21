package handler

import (
	"github.com/gorilla/websocket"
	"github.com/joeyCheek888/go-std/httpx/websocket/manager"
	"github.com/joeyCheek888/go-std/log"
	"net/http"
	"time"
)

func Upgrade(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			log.Logger.Info("升级协议",
				log.String("url", req.URL.String()),
				log.String("host", req.Host),
				log.Any("ua", r.Header["User-Agent"]),
				log.Any("referer", r.Header["Referer"]))
			return true
		},
		EnableCompression: false,
	}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	token := req.PathValue("token")
	if token == "" {
		http.NotFound(w, req)
		return
	}
	log.Logger.Info("webSocket 建立连接:", log.String("addr", conn.RemoteAddr().String()))
	currentTime := time.Now().Unix()
	c := manager.NewMember(conn.RemoteAddr().String(), token, conn, currentTime)
	go c.Read()
	go c.Write()

	// 用户连接事件
	manager.Manager.Register <- c
}
