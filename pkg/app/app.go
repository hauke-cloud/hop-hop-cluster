package app

import (
	"context"

	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	http "github.com/hauke-cloud/hop-hop-cluster/pkg/api"
	logger "github.com/hauke-cloud/hop-hop-cluster/pkg/logger"
	services "github.com/hauke-cloud/hop-hop-cluster/pkg/usecase/interface"
)

type StartApp struct {
	httpServer     *http.ServerHTTP
	clusterUseCase services.ClusterUseCase
	logger         logger.Logger
	config         config.Config
}

func NewStartApp(httpServer *http.ServerHTTP, clusterUseCase services.ClusterUseCase, config config.Config, logger logger.Logger) *StartApp {
	return &StartApp{
		httpServer:     httpServer,
		clusterUseCase: clusterUseCase,
		config:         config,
		logger:         logger,
	}
}

func (a *StartApp) Initialize() error {
	a.logger.InitLogger()

	a.logger.Infof("Starting to watch cluster apis")

	// Start API watching
	err := a.clusterUseCase.Initialize()

	return err
}

func (a *StartApp) Run(ctx context.Context) {
	a.httpServer.Start()
	a.clusterUseCase.Start()
}

func (a *StartApp) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down the http server")
	if err := a.httpServer.Stop(ctx); err != nil {
		a.logger.Fatalf("Failed to stop the http server: %v", err)
		return err
	}

	return nil
}
