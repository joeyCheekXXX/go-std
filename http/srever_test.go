package http

import "testing"

func TestServer(t *testing.T) {

	conf := &Config{
		Addr:             "",
		Port:             "",
		RouterPrefix:     "",
		Mode:             "",
		LimitCountIP:     0,
		LimitTimeIP:      0,
		UseMultipoint:    false,
		EnableSwaggerDoc: false,
	}

	s := NewServer(conf)

	s.Start()
}
