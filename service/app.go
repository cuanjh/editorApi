package service

import (
	"github.com/facebookgo/inject"
)

//s依赖注入
func AppService() *Object {
	var g inject.Graph
	baseService := BaseService{}
	g.Provide(&inject.Object{Value: baseService})
	g.Populate()
	return &Object{
		AppService: &baseService,
	}
}

type Object struct {
	AppService *BaseService
}
