package groupx

import (
	"github.com/joeyCheek888/go-std.git/log"
	"github.com/zeromicro/go-zero/core/service"
)

type Service service.Service

type Group struct {
	Services []service.Service
}

func NewServiceGroup() *Group {
	return &Group{
		Services: make([]service.Service, 0),
	}
}

func (g Group) Start() {
	log.Logger.Info("service group start", log.Int("count", len(g.Services)))
	for _, _service := range g.Services {
		_service.Start()
	}
}

func (g Group) Stop() {
	log.Logger.Info("service group stop", log.Int("count", len(g.Services)))
	for _, _service := range g.Services {
		_service.Stop()
	}
}

func (g Group) Add(services ...Service) {
	for _, item := range services {
		g.Services = append(g.Services, item)
	}
}
