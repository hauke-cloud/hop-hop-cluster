package usecase

import (
	"context"
	"sync"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/hauke-cloud/hop-hop-cluster/internal/config"
	domain "github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
	logger "github.com/hauke-cloud/hop-hop-cluster/pkg/logger"
	interfaces "github.com/hauke-cloud/hop-hop-cluster/pkg/repository/interface"
	services "github.com/hauke-cloud/hop-hop-cluster/pkg/usecase/interface"
)

type ClusterUseCase struct {
	clusterRepo interfaces.ClusterRepository
	validator   *validator.Validate
	logger      logger.Logger
	mu          sync.Mutex
	shutdown    chan struct{}
	done        chan struct{}
	members     []config.Member
	config      config.Config
}

func NewClusterUseCase(clusterRepo interfaces.ClusterRepository, logger logger.Logger, config config.Config) services.ClusterUseCase {
	validator := validator.New()

	return &ClusterUseCase{
		clusterRepo: clusterRepo,
		validator:   validator,
		logger:      logger,
		shutdown:    make(chan struct{}),
		done:        make(chan struct{}),
		config:      config,
	}
}

func (c *ClusterUseCase) Initialize() error {
	c.logger.Debugf("Initializing cluster use case")

	c.members = c.config.Cluster.Members
	return nil
}

func (c *ClusterUseCase) Start() error {
	c.logger.Debugf("Starting cluster use case")
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger.Debugf("Fetching members from config")
	for _, member := range c.members {
		c.logger.Debugf("Starting to watch cluster api of host: %s", member.IP)
		go c.runWatchAPI(member, time.Duration(c.config.Cluster.Interval)*time.Second, c.config.Cluster.Retries)
	}
	<-c.done
	return nil
}

func (c *ClusterUseCase) runWatchAPI(member config.Member, interval time.Duration, max_retries int) {
	retries := max_retries

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	defer close(c.done)
	for retries > 0 {
		select {
		case <-c.shutdown:
			c.logger.Infof("Finished watching cluster api of host: %s", member.IP)
			return
		default:
			c.logger.Debugf("Fetching cluster api (host: %s, retries: %d)", member.IP, retries)
			err := c.clusterRepo.GetAndSaveAPI(ctx, member, c.config.General.ListenPort)
			if err != nil {
				if retries > 0 {
					c.logger.Debugf("Error fetching cluster api (host: %s, retries: %d, message: %s). Retrying...", member.IP, retries, err)
					retries--
				} else {
					c.logger.Errorf("Error fetching cluster api (host: %s, retries: %d, message: %s). End of retried reached. Stopping...", member.IP, retries, err)
					continue
				}
			}
		}
		time.Sleep(interval)
	}
}

func (c *ClusterUseCase) Stop() {
	close(c.shutdown)
	<-c.done
}

func (c *ClusterUseCase) GetThis(ctx context.Context) domain.Cluster {
	cluster := c.clusterRepo.GetThis(ctx)
	return cluster
}

func (c *ClusterUseCase) FindAll(ctx context.Context) ([]domain.Cluster, error) {
	clusters, err := c.clusterRepo.FindAll(ctx)
	return clusters, err
}

func (c *ClusterUseCase) FindByName(ctx context.Context, name string) (domain.Cluster, error) {
	cluster, err := c.clusterRepo.FindByName(ctx, name)
	return cluster, err
}

func (c *ClusterUseCase) Save(ctx context.Context, cluster domain.Cluster) (domain.Cluster, error) {
	cluster, err := c.clusterRepo.Save(ctx, cluster)

	return cluster, err
}

func (c *ClusterUseCase) Delete(ctx context.Context, cluster domain.Cluster) error {
	err := c.clusterRepo.Delete(ctx, cluster)

	return err
}
