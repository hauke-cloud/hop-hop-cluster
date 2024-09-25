package interfaces

import (
	"context"

	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	"github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
)

type ClusterRepository interface {
	GetThis(ctx context.Context) domain.Cluster
	FindAll(ctx context.Context) ([]domain.Cluster, error)
	FindByName(ctx context.Context, name string) (domain.Cluster, error)
	Save(ctx context.Context, cluster domain.Cluster) (domain.Cluster, error)
	Delete(ctx context.Context, cluster domain.Cluster) error
	GetAndSaveAPI(ctx context.Context, member config.Member, port int) error
}
