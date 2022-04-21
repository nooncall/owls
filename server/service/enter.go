package service

import (
	"github.com/qingfeng777/owls/server/service/autocode"
	"github.com/qingfeng777/owls/server/service/example"
	"github.com/qingfeng777/owls/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	AutoCodeServiceGroup autocode.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
