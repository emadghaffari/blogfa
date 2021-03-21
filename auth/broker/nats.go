package broker

import (
	"blogfa/auth/config"
	"strings"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
)

var (
	Nats  NatsBroker = &nts{}
	oncen sync.Once
)

type NatsBroker interface{}

type nts struct {
	conn *nats.Conn
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
			return
		}

	})

	return err
}

func (n *nts) GetConnection() *nats.Conn {
	return n.conn
}
