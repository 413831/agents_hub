package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainTask "proyecto_ia/internal/domain/task"
)

// EnqueueTaskUseCase coloca una tarea en el broker para procesamiento asíncrono
type EnqueueTaskUseCase struct {
	broker domainTask.TaskBroker
	repo   domainTask.JobRepository
}

func NewEnqueueTaskUseCase(broker domainTask.TaskBroker, repo domainTask.JobRepository) *EnqueueTaskUseCase {
	return &EnqueueTaskUseCase{
		broker: broker,
		repo:   repo,
	}
}

func (uc *EnqueueTaskUseCase) Execute(ctx context.Context, taskType domainTask.TaskType, payload interface{}) (*domainTask.Job, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize task payload: %w", err)
	}

	job := domainTask.Job{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Type:      taskType,
		Status:    domainTask.StatusPending,
		Progress:  0,
		Payload:   bytes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if uc.repo != nil {
		if err := uc.repo.Save(ctx, job); err != nil {
			return nil, fmt.Errorf("failed to save job to repository: %w", err)
		}
	}

	// Publicar en el topic correspondiente
	topic := fmt.Sprintf("tasks.%s", taskType)
	if err := uc.broker.Publish(ctx, topic, job); err != nil {
		return nil, fmt.Errorf("failed to publish job to broker: %w", err)
	}

	return &job, nil
}
