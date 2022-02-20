//go:build wireinject
// +build wireinject

package main

import (
	"github.com/cycade/go-samples/internal/biz"
	"github.com/cycade/go-samples/internal/data"
	"github.com/cycade/go-samples/internal/service"
	"github.com/google/wire"
)

func initApp() service.DeploymentService {
	wire.Build(data.NewDeploymentRepo, biz.NewDeploymentUsecase, service.NewDeploymentService)
	return service.DeploymentService{}
}
