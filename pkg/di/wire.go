//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	config "github.com/hauke-cloud/hop-hop-cluster/internal/config"
	http "github.com/hauke-cloud/hop-hop-cluster/pkg/api"
	handler "github.com/hauke-cloud/hop-hop-cluster/pkg/api/handler"
	app "github.com/hauke-cloud/hop-hop-cluster/pkg/app"
	"github.com/hauke-cloud/hop-hop-cluster/pkg/db"
	logger "github.com/hauke-cloud/hop-hop-cluster/pkg/logger"
	repository "github.com/hauke-cloud/hop-hop-cluster/pkg/repository"
	usecase "github.com/hauke-cloud/hop-hop-cluster/pkg/usecase"
)

func Initialize(cfg config.Config) (*app.StartApp, error) {
	wire.Build(db.ConnectDatabase,
		// Logger
		logger.NewAppLogger,
		wire.Bind(new(logger.Logger), new(*logger.AppLogger)),
		// Cluster
		repository.NewClusterRepository,
		usecase.NewClusterUseCase,
		handler.NewClusterHandler,
		// HTTP server
		http.NewServerHTTP,
		// App
		app.NewStartApp,
	)

	return &app.StartApp{}, nil
}
