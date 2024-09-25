package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/hauke-cloud/hop-hop-cluster/cmd/api/docs"
	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	handler "github.com/hauke-cloud/hop-hop-cluster/pkg/api/handler"
)

type ServerHTTP struct {
	engine *http.Server
}

func NewServerHTTP(clusterHandler *handler.ClusterHandler, config config.Config) *ServerHTTP {
	router := gin.Default()
	router.Use(gin.Logger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Simple ping
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := router.Group("/api")

	// cluster routes
	api.GET("cluster", clusterHandler.GetThis)
	api.GET("clusters", clusterHandler.GetCluster)

	engine := &http.Server{
		Addr:      fmt.Sprintf(":%d", config.General.ListenPort),
		Handler:   router,
		TLSConfig: config.ClientTLS,
	}

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	go func() {
		if err := sh.engine.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (sh *ServerHTTP) Stop(ctx context.Context) error {
	return sh.engine.Shutdown(ctx)
}
