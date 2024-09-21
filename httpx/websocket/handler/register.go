package handler

import (
	"github.com/joeyCheek888/go-std/httpx/websocket/client"
	"github.com/joeyCheek888/go-std/log"
)

type registerCallbackFun func(client *client.Client)

var RegisterCallback = func(client *client.Client) {
	log.Logger.Info("用户连接", log.String("addr", client.Addr))
}

func SetRegisterCallback(callback registerCallbackFun) {
	RegisterCallback = callback
}
