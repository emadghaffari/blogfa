package consul

import (
	"blogfa/auth/config"
	zapLogger "blogfa/auth/pkg/logger"
	"sync"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var (
	Storage store = &consul{}
	once    sync.Once
)

// store interface is interface for store things into mysql
type store interface {
	Connect(config config.GlobalConfig) error
}

// mysql struct
type consul struct {
	client *api.Client
}

func (c *consul) Connect(config config.GlobalConfig) error {
	var err error

	once.Do(func() {
		// configs to connect to consul
		c.client, err = api.NewClient(&api.Config{
			Address: config.Consul.Address,
			Scheme:  config.Consul.Scheme,
		})

		if err != nil {
			logger := zapLogger.GetZapLogger(config.Debug())
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())

			return
		}

		if err := c.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			Address: config.Consul.Address,
			ID:      "token.Generate(10)", // Unique for each node
			Name:    config.Service.Name,  // Can be service type
			Tags:    config.Consul.Tags,
			Check: &api.AgentServiceCheck{
				HTTP:     config.Consul.CheckHttp,
				Interval: config.Consul.CheckInterval,
			},
		}); err != nil {
			logger := zapLogger.GetZapLogger(config.Debug())
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
			return
		}

	})

	return err
}

func (c *consul) GetClient() *api.Client {
	return c.client
}

func (c *consul) CreateSesstion(a api.SessionEntry) (string, *api.WriteMeta, error) {
	sessionID, meta, err := c.client.Session().Create(&a, nil)
	if err != nil {
		logger := zapLogger.GetZapLogger(config.Global.Debug())
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return "", nil, err
	}

	return sessionID, meta, err
}
