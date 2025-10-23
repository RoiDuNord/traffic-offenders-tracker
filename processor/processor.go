package processor

import (
	"context"
	"fmt"
	"log/slog"
	"speed_violation_tracker/broker"
	"speed_violation_tracker/cat"
	wp "speed_violation_tracker/workerPool"
	"time"
)

type Message = cat.Message

type DB interface {
	InsertMessage(key string, value []byte) (int, error)
	Close() error
}

type BrokerChan = broker.BrokerChan

type Broker interface {
	Subscribe() (BrokerChan, error)
	Close() error
}

type Pool interface {
	Create()
	Handle(msg cat.Message, cancel context.CancelFunc)
	Wait(ctx context.Context)
	Stats()
}

type Processor struct {
	ctx    context.Context
	cancel context.CancelFunc

	db     DB
	broker Broker
	wp     Pool

	maxOffenders     int
	offendersChan    chan struct{}
	offendersHandled int

	closeTimeout int // seconds
}

func New(ctx context.Context, db DB, broker Broker, maxOffenders, closeTimeout int) *Processor {
	ctx, cancel := context.WithCancel(ctx)
	return &Processor{
		ctx:    ctx,
		cancel: cancel,
		db:     db,
		broker: broker,

		maxOffenders:  maxOffenders,
		offendersChan: make(chan struct{}),
		closeTimeout:  closeTimeout,
	}
}

func (p *Processor) Run() error {
	messagesChan, err := p.broker.Subscribe()
	if err != nil {
		return err
	}

	p.wp = wp.New(p.processMessage)
	p.wp.Create()

	go func() {
		for msg := range messagesChan {
			p.wp.Handle(msg, p.cancel)
		}
		// если брокер закрыл канал, то закрываем обработку
		p.cancel()
	}()

	go func() {
		for range p.offendersChan {
			p.offendersHandled++
			if p.maxOffenders <= 0 {
				continue
			}
			if p.offendersHandled >= p.maxOffenders {
				slog.Info("Max offenders reached, program will close.")
				p.cancel()
			}
		}
	}()

	workerPoolDone := make(chan struct{})

	go func() {
		<-p.ctx.Done()
		slog.Info(fmt.Sprintf("Context done %s", p.ctx.Err().Error()))
		p.broker.Close()

		slog.Info("Wait for worker pool ended: 10s timeout")
		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.closeTimeout))
		defer cancel()
		p.wp.Wait(closeCtx)

		workerPoolDone <- struct{}{}
	}()

	<-workerPoolDone
	close(p.offendersChan)

	p.wp.Stats()

	return nil
}
