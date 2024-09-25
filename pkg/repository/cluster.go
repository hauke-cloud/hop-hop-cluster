package repository

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	domain "github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
	interfaces "github.com/hauke-cloud/hop-hop-cluster/pkg/repository/interface"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type clusterDatabase struct {
	DB     *gorm.DB
	config config.Config
	This   domain.Cluster
	client *resty.Client
}

func NewClusterRepository(DB *gorm.DB, config config.Config) interfaces.ClusterRepository {
	return &clusterDatabase{
		DB,
		config,
		domain.Cluster{
			Name:      config.General.NodeName,
			IPAddress: config.General.NodeIP,
			Priority:  rand.Intn(10000),
			Status:    domain.Waiting,
		},
		resty.New().SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: config.ClientTLS.Certificates,
				RootCAs:      config.ClientTLS.RootCAs,
			},
		}),
	}
}

func (c *clusterDatabase) GetAndSaveAPI(ctx context.Context, member config.Member, port int) error {
	c.client.R().SetDebug(true)
	resp, err := c.client.R().
		SetResult(&domain.Cluster{}).
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("https://%s:%d/api/cluster", member.IP, port))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New("failed to fetch setting: " + resp.Status())
	}

	_, err = c.Save(ctx, *resp.Result().(*domain.Cluster))
	if err != nil {
		return err
	}

	return nil
}

func (c *clusterDatabase) GetThis(ctx context.Context) domain.Cluster {
	return c.This
}

func (c *clusterDatabase) FindAll(ctx context.Context) ([]domain.Cluster, error) {
	var clusters []domain.Cluster

	err := c.DB.Find(&clusters).Error

	return clusters, err
}

func (c *clusterDatabase) FindByName(ctx context.Context, name string) (domain.Cluster, error) {
	var cluster domain.Cluster

	err := c.DB.Where("name = ?", name).First(&cluster).Error

	return cluster, err
}

func (c *clusterDatabase) Save(ctx context.Context, cluster domain.Cluster) (domain.Cluster, error) {
	err := c.DB.Where("name=?", cluster.Name).Save(&cluster).Error

	return cluster, err
}

func (c *clusterDatabase) Delete(ctx context.Context, cluster domain.Cluster) error {
	err := c.DB.Delete(&cluster).Error

	return err
}
