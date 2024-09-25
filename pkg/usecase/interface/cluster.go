package interfaces

import (
	"context"

	domain "github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
)

type ClusterUseCase interface {
	GetThis(ctx context.Context) domain.Cluster
	FindAll(ctx context.Context) ([]domain.Cluster, error)
	FindByName(ctx context.Context, name string) (domain.Cluster, error)
	FindClusterLeader(ctx context.Context) (domain.Cluster, error)
	GetClusterStatus(ctx context.Context) (domain.Status, error)
	Save(ctx context.Context, cluster domain.Cluster) (domain.Cluster, error)
	Delete(ctx context.Context, cluster domain.Cluster) error
	Initialize() error
	Start() error
}
