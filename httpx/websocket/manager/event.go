package manager

import (
	"github.com/joeyCheek888/go-std/log"
	"net/http"
)

// Event
// @Description: manager对客户端事件的处理结构定义
type Event struct {
	RegisterChan       chan *Member                         // 连接连接处理 chan
	UnregisterChan     chan *Member                         // 断开连接处理 chan
	BroadcastChan      chan []byte                          // 广播 向全部成员发送数据 chan
	registerConnFunc   func(member *Member)                 // 注册连接触发的方法
	closeConnFunc      func(member *Member)                 // 关闭连接触发的方法
	processMessageFunc func(member *Member, message []byte) // 客户端发送消息时的方法
	CheckTokenFunc     func(token string) bool              // 尝试注册连接前校验检查token的方法
	CheckOriginFunc    func(r *http.Request) bool
}

func NewEvent() *Event {
	return &Event{
		RegisterChan:   make(chan *Member, 1000),
		UnregisterChan: make(chan *Member, 1000),
		BroadcastChan:  make(chan []byte, 1000),
		registerConnFunc: func(member *Member) {
			log.Logger.Info("用户连接", log.String("addr", member.Addr), log.Int64("ID", member.ID), log.String("token", member.Token))
		},
		closeConnFunc: func(member *Member) {
			log.Logger.Info("关闭连接", log.String("addr", member.Addr))
		},
		processMessageFunc: func(member *Member, message []byte) {
			log.Logger.Info("收到消息", log.Int64("ID", member.ID), log.String("addr", member.Addr), log.String("message", string(message)))
		},
		CheckTokenFunc: func(token string) bool {
			if token == "" {
				return false
			}
			return true
		},
		CheckOriginFunc: func(r *http.Request) bool {
			log.Logger.Info("请求地址", log.String("host", r.Host), log.String("url", r.URL.String()), log.Any("ua", r.Header["User-Agent"]), log.Any("referer", r.Header["Referer"]))
			return true
		},
	}
}
