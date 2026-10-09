package broker

import (
	"context"
	"fmt"
	"sync"

	domainTask "proyecto_ia/internal/domain/task"
)

// InProcessBroker es un Message Broker de alto rendimiento en memoria/canales concurrentes
// Diseñado como broker embebido de cero configuración para el monolito de escritorio.
// Cumple al 100% el contrato de TaskBroker al igual que NATS o Redis.
type InProcessBroker struct {
	mu          sync.RWMutex
	subscribers map[string][]func(job domainTask.Job)
	jobQueue    chan struct {
		topic string
		job   domainTask.Job
	}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewInProcessBroker(bufferSize int) *InProcessBroker {
	ctx, cancel := context.WithCancel(context.Background())
	b := &InProcessBroker{
		subscribers: make(map[string][]func(job domainTask.Job)),
		jobQueue: make(chan struct {
			topic string
			job   domainTask.Job
		}, bufferSize),
		ctx:    ctx,
		cancel: cancel,
	}

	go b.dispatcher()
	return b
}

func (b *InProcessBroker) dispatcher() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case item, ok := <-b.jobQueue:
			if !ok {
				return
			}
			b.mu.RLock()
			handlers, exists := b.subscribers[item.topic]
			b.mu.RUnlock()

			if exists {
				for _, h := range handlers {
					go h(item.job) // Despachar en goroutine
				}
			}
		}
	}
}

func (b *InProcessBroker) Publish(ctx context.Context, topic string, job domainTask.Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.ctx.Done():
		return fmt.Errorf("broker is closed")
	case b.jobQueue <- struct {
		topic string
		job   domainTask.Job
	}{topic: topic, job: job}:
		return nil
	}
}

func (b *InProcessBroker) Subscribe(ctx context.Context, topic string, handler func(job domainTask.Job)) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}

func (b *InProcessBroker) Close() error {
	b.cancel()
	close(b.jobQueue)
	return nil
}
