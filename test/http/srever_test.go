package http

import (
	"github.com/joeyCheek888/go-std/httpx"
	"testing"
)

func TestServer(t *testing.T) {

	conf := &httpx.Config{
		Addr:             "",
		Port:             "",
		RouterPrefix:     "",
		Mode:             "",
		EnableSwaggerDoc: false,
	}

	s := httpx.NewServer(conf)

	s.Start()
}
