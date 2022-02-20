package data

import (
	"github.com/cycade/go-samples/internal/biz"
	"github.com/google/wire"
)

var Providers = wire.NewSet(NewDeploymentRepo)

func NewDeploymentRepo() biz.DeploymentRepo {
	return deploymentRepo{}
}

type deploymentRepo struct {
}

func (r deploymentRepo) Create(d *biz.Deployment) (string, error) {
	return "42", nil
}

func (r deploymentRepo) GetStatus(id string) (biz.Deployment, error) {
	return biz.DefaultDeployment, nil
}

func (r deploymentRepo) Update(d *biz.Deployment) error {
	return nil
}

var _ biz.DeploymentRepo = deploymentRepo{}
