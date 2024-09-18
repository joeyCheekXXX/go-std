package websocket

import (
	"fmt"
	"github.com/goccy/go-json"
	"go-std/http/websocket/types"
	"time"
)

// defaultProcessData 处理数据
func defaultProcessData(client *Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("websocket defaultProcessData recover 处理数据 ", r)
		}
	}()
	request := &types.Request{}
	if err := json.Unmarshal(message, request); err != nil {
		client.SendMsg([]byte("数据不合法"))
		return
	}

	cmd := request.Cmd
	data := request.Data

	if cmd == "heartbeat" {
		client.Heartbeat(time.Now().Unix())
		heartbeatCallback(client)
		return
	}

	processMessageCallback(client, cmd, data)

	return
}
