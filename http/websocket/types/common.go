package types

import "github.com/goccy/go-json"

type Request struct {
	Cmd  string      `json:"cmd" example:"subscribe"` // 客户端请求的动作 如 login(登录)/heartbeat(心跳)/subscribe(订阅)/unsubscribe(取消订阅)
	Data interface{} `json:"data,omitempty"`          // 数据 json
}

// Response 响应数据体
type Response struct {
	Cmd     string      `json:"cmd"` // 客户端请求的动作 如 login(登录)/heartbeat(心跳)/subscribe(订阅)/unsubscribe(取消订阅)
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"` // 数据 json
}

func (r Response) Marshal() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}

// NewResponse 创建新的响应
func NewResponse(cmd string, success bool, msg string, data interface{}) *Response {

	uint8Data, ok := data.([]uint8)
	if ok && uint8Data == nil {
		data = make(map[string]any)
	}

	return &Response{
		Cmd:     cmd,
		Success: success,
		Msg:     msg,
		Data:    data,
	}
}
