package broker

import (
	"blogfa/auth/config"
	zapLogger "blogfa/auth/pkg/logger"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	Nats  NatsBroker = &nts{}
	oncen sync.Once
)

type NatsBroker interface {
	Connect() error
	Conn() *nats.Conn
}

type nts struct {
	conn *nats.Conn
	ec   *nats.EncodedConn
}

func (n *nts) Connect() error {
	var err error
	oncen.Do(func() {
		opts := nats.Options{
			Name:         config.Global.Service.Name,
			Secure:       config.Global.Nats.Auth,
			User:         config.Global.Nats.Username,
			Password:     config.Global.Nats.Password,
			MaxReconnect: 10,
			Url:          strings.Join(config.Global.Nats.Endpoints, ","),
			PingInterval: time.Minute * 10,
		}

		n.conn, err = opts.Connect()
		if err != nil {
			logger := zapLogger.GetZapLogger(false)
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
			return
		}

	})

	return err
}

func (n *nts) EncodedConn() error {
	if n.conn == nil {
		logger := zapLogger.GetZapLogger(false)
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit("ERROR")
		return fmt.Errorf("nats not connected, connect first")
	}

	var err error
	n.ec, err = nats.NewEncodedConn(n.conn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	return nil
}

func (n *nts) Conn() *nats.Conn {
	return n.conn
}
func (n *nts) ECConn() *nats.EncodedConn {
	return n.ec
}

func (n *nts) Publish(ctx context.Context, subject string, value interface{}) error {
	if err := n.ec.Publish(subject, &value); err != nil {
		logger := zapLogger.GetZapLogger(false)
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}

	return nil
}
