package handler

import (
	"github.com/joeyCheek888/go-std/httpx/websocket/client"
	"github.com/joeyCheek888/go-std/log"
)

type closeConnCallback func(client *client.Client)

var CloseConnCallback = func(client *client.Client) {
	log.Logger.Info("关闭连接", log.String("addr", client.Addr))
}

func SetCloseConnCallback(callback closeConnCallback) {
	CloseConnCallback = callback
}
