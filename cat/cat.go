// Package cat simulates the operation of a message broker.
// Minimal functionality implemented.
package cat

import "sync"

type Cat struct {
	connected bool
	dataCh    CatChan
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

type Message struct {
	b []byte
}

func (m *Message) Ack() error {
	return nil
}

type CatChan chan Message

// New returns new cat-broker entity
func New() *Cat {
	return new(Cat)
}

// Connect connects to server using conn
func (c *Cat) Connect(conn string) error {
	c.connected = true
	return nil
}

// Subscribe returns CatChan and start broadcast
func (c *Cat) Subscribe() (CatChan, error) {
	if !c.connected {
		return nil, ErrHasNoConn
	}

	c.dataCh = make(CatChan)
	c.stopCh = make(chan struct{})

	c.wg.Add(1)
	go c.broadcast()

	return c.dataCh, nil
}

// Close closes cat connection
func (c *Cat) Close() error {
	if !c.connected {
		return nil
	}
	c.connected = !c.connected

	close(c.stopCh)
	c.wg.Wait()

	c.dataCh = nil
	return nil
}

// Bytes returns a bytes body message
func (m *Message) Bytes() []byte {
	return m.b
}
