package handler

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/joeyCheek888/go-std/httpx/websocket/manager"
	"github.com/joeyCheek888/go-std/httpx/websocket/types"
	"github.com/joeyCheek888/go-std/log"
)

// DefaultProcessData 处理数据
func DefaultProcessData(client *manager.Member, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("websocket recover 处理数据 ", r)
		}
	}()
	request := &types.Request{}
	if err := json.Unmarshal(message, request); err != nil {
		client.SendMsg([]byte("数据不合法"))
		return
	}

	cmd := request.Cmd
	data := request.Data

	ProcessMessageCallback(client, cmd, data)

	return
}

type processMessageCallbackFun func(client *manager.Member, cmd string, data any)

var ProcessMessageCallback = func(client *manager.Member, cmd string, data any) {
	log.Logger.Info(
		"收到消息",
		log.String("addr", client.Addr),
		log.Any("ID", client.ID),
		log.String("cmd", cmd), log.Any("data", data),
	)
}

func SetProcessMessageCallback(callback processMessageCallbackFun) {
	ProcessMessageCallback = callback
}
