package service

import (
	"github.com/cycade/go-samples/api"
	"github.com/cycade/go-samples/internal/biz"
	"github.com/google/wire"
)

var Providers = wire.NewSet(NewDeploymentService)

func NewDeploymentService(usecase biz.DeploymentUsecase) DeploymentService {
	return DeploymentService{usecase: usecase}
}

type DeploymentService struct {
	usecase biz.DeploymentUsecase
}

func (s *DeploymentService) Restart(id string) *api.DeploymentRestartResp {
	err := s.usecase.Restart(id)
	if err != nil {
		return &api.DeploymentRestartResp{Status: "failed"}
	}

	return &api.DeploymentRestartResp{Status: "success"}
}
