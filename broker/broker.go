package broker

import (
	"fmt"
	"speed_violation_tracker/cat"
)

type BrokerChan = cat.CatChan

type Broker interface {
	Connect(conn string) error
	Subscribe() (BrokerChan, error)
	Close() error
}

type BrokerService struct {
	broker Broker
}

func New(broker Broker) *BrokerService {
	return &BrokerService{broker: broker}
}

func (bm *BrokerService) Subscribe() (BrokerChan, error) {
	ch, err := bm.broker.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("BrokerService: subscribe error: %w", err)
	}
	return ch, nil
}

func (bm *BrokerService) Close() error {
	return bm.broker.Close()
}
