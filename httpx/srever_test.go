package httpx

import "testing"

func TestServer(t *testing.T) {

	conf := &Config{
		Addr:             "",
		Port:             "",
		RouterPrefix:     "",
		Mode:             "",
		EnableSwaggerDoc: false,
	}

	s := NewServer(conf)

	s.Start()
}
