package biz

import (
	"fmt"
	"strings"

	"github.com/google/wire"
)

type Deployment struct {
	Id          string
	Image       string
	Replicas    int32
	StartScript string
	ArtifactId  string
}

var DefaultDeployment = Deployment{
	Id:          "93",
	Image:       "alpine",
	Replicas:    3,
	StartScript: "sleep 1000000",
	ArtifactId:  "1090",
}

type DeploymentRepo interface {
	Create(d *Deployment) (string, error)
	GetStatus(id string) (Deployment, error)
	Update(d *Deployment) error
}

var images = map[string]string{
	"java":   "openjdk:18-oracle",
	"nodejs": "nodejs:16-bullseye",
}

var Providers = wire.NewSet(NewDeploymentUsecase)

func NewDeploymentUsecase(repo DeploymentRepo) DeploymentUsecase {
	return DeploymentUsecase{repo: repo}
}

type DeploymentUsecase struct {
	repo DeploymentRepo
}

func (du *DeploymentUsecase) UpdateKind(id string, kind string) error {
	repo, err := du.repo.GetStatus(id)
	if err != nil {
		return err
	}

	image, ok := images[kind]
	if ok {
		repo.Image = image
	} else {
		return fmt.Errorf("未能找到 kind %s 对应的镜像", kind)
	}

	return du.repo.Update(&repo)
}

// 更改容器启动命令导致 yaml 变化，pod 会自动重启
func (du *DeploymentUsecase) Restart(id string) error {
	repo, err := du.repo.GetStatus(id)
	if err != nil {
		return err
	}

	if strings.HasSuffix(repo.StartScript, " ") {
		length := len(repo.StartScript)
		repo.StartScript = repo.StartScript[0 : length-1]
	} else {
		repo.StartScript = repo.StartScript + " "
	}

	return du.repo.Update(&repo)
}
