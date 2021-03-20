package broker

import "sync"

var (
	Nats NatsBroker = &nats{}
	nts  sync.Once
)

type NatsBroker interface{}

type nats struct{}

func (n *nats) Connect() {
	nts.Do(func() {})
}
