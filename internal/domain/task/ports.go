package task

import "context"

// HandlerFunc callback que procesa un trabajo
type HandlerFunc func(ctx context.Context, job Job) (result []byte, err error)

// TaskBroker puerto para el bus de mensajes (NATS o Redis)
type TaskBroker interface {
	Publish(ctx context.Context, topic string, job Job) error
	Subscribe(ctx context.Context, topic string, handler func(job Job)) error
	Close() error
}

// JobRepository persistencia del estado de los trabajos
type JobRepository interface {
	Save(ctx context.Context, job Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	ListRecent(ctx context.Context, limit int) ([]Job, error)
	UpdateStatus(ctx context.Context, id string, status TaskStatus, progress int, result []byte, errStr string) error
}
