package wp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"speed_violation_tracker/cat"
	"speed_violation_tracker/models"
)

type Worker struct {
	id                int
	vehiclesCompleted int
}

var workers = []*Worker{
	{id: 1},
	{id: 2},
	{id: 3},
	{id: 4},
}

type Pool struct {
	pool    chan *Worker
	handler func(workerID int, data cat.Message) error
}

func New(handler func(int, cat.Message) error) *Pool {
	return &Pool{
		handler: handler,
		pool:    make(chan *Worker, len(workers)),
	}
}

func (p *Pool) Create() {
	for _, w := range workers {
		p.pool <- w
	}
}

func (p *Pool) Handle(msg cat.Message, cancel context.CancelFunc) {
	w := <-p.pool

	go func() {
		defer func() {
			p.pool <- w
		}()

		if err := p.handler(w.id, msg); err != nil {
			var fatalErr *models.FatalError
			if errors.As(err, &fatalErr) {
				slog.Error(fmt.Sprintf("Handle: fatal error occurred: %s", err.Error()))
				cancel()
			} else {
				slog.Error(fmt.Sprintf("Handle: error occurred: %s", err.Error()))
			}
			return
		}

		w.vehiclesCompleted++
	}()
}

func (p *Pool) Wait(ctx context.Context) {
	closeCh := make(chan struct{})
	go func() {
		for range len(workers) {
			<-p.pool
		}
		closeCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-closeCh:
		return
	}
}

func (p *Pool) Stats() {
	slog.Info("__________Results__________")
	for _, w := range workers {
		slog.Info(fmt.Sprintf("Worker %d processed %d vehicles\n", w.id, w.vehiclesCompleted))
	}
}
