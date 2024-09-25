package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	"github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
	service "github.com/hauke-cloud/hop-hop-cluster/pkg/usecase/interface"
)

type ClusterHandler struct {
	clusterUseCase service.ClusterUseCase
	config         config.Config
}

func NewClusterHandler(clusterUseCase service.ClusterUseCase, config config.Config) *ClusterHandler {
	return &ClusterHandler{
		clusterUseCase: clusterUseCase,
		config:         config,
	}
}

type ClusterResponse struct {
	Name      string        `json:"name" copier:"must"`
	IPAddress string        `json:"ip_address" copier:"must"`
	Priority  int           `json:"priority" copier:"must"`
	Status    domain.Status `json:"status" copier:"must"`
}

// FindAll godoc
// @summary Get all clusters
// @description Get all clusters
// @tags clusters
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/clusters [get]
// @response 200 {object} []TeamplteResponse "OK"
func (cr *ClusterHandler) GetCluster(c *gin.Context) {
	clusters, err := cr.clusterUseCase.FindAll(c)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldnt fetch cluster", "message": err})
	} else {
		response := []ClusterResponse{}
		copier.Copy(&response, &clusters)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *ClusterHandler) GetThis(c *gin.Context) {
	cluster := cr.clusterUseCase.GetThis(c)
	response := ClusterResponse{}
	copier.Copy(&response, &cluster)
	c.JSON(http.StatusOK, response)
}
