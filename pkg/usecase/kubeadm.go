package usecase

import (
	services "github.com/hauke-cloud/hop-hop-cluster/pkg/usecase/interface"
)

type KubeadmUseCase struct{}

func NewKubeadmUseCase() services.KubeadmUseCase {
	return &KubeadmUseCase{}
}
