package service

import (
	"github.com/nooncall/owls/go/service/autocode"
	"github.com/nooncall/owls/go/service/example"
	"github.com/nooncall/owls/go/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	AutoCodeServiceGroup autocode.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
