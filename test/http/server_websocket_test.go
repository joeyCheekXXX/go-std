package http

import (
	"github.com/joeyCheek888/go-std/httpx"
	"testing"
)

func TestWebsocket(t *testing.T) {

	conf := &httpx.Config{
		Addr:             "",
		Port:             "",
		RouterPrefix:     "",
		Mode:             "",
		EnableSwaggerDoc: false,
	}

	s := httpx.NewServer(conf)
	defer s.Stop()

	s.EnableWebsocket("/test/ws/:token")

	s.WebsocketManger().WithCheckTokenFunc(func(token string) bool {
		if token == "" || token == "123123" {
			return false
		}

		return true
	})

	//s.WebsocketManger().WithRegisterConnFunc(func(member *manager.Member) {
	//
	//})
	//
	//s.WebsocketManger().WithCloseConnFunc(func(member *manager.Member) {
	//
	//})
	//
	//s.WebsocketManger().WithProcessMessageFunc(func(member *manager.Member, message []byte) {
	//
	//})

	s.Start()

}
