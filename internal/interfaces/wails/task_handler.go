package wails

import (
	"context"

	domainTask "proyecto_ia/internal/domain/task"
	usecaseTask "proyecto_ia/internal/usecase/task"
)

// TaskHandler expone la gestión de tareas asíncronas y Message Broker al frontend
type TaskHandler struct {
	enqueueUseCase *usecaseTask.EnqueueTaskUseCase
	repo           domainTask.JobRepository
	app            *App
}

func NewTaskHandler(
	enqueueUseCase *usecaseTask.EnqueueTaskUseCase,
	repo domainTask.JobRepository,
	app *App,
) *TaskHandler {
	return &TaskHandler{
		enqueueUseCase: enqueueUseCase,
		repo:           repo,
		app:            app,
	}
}

// EnqueueLLMTask encola una tarea de inferencia pesada en segundo plano
func (h *TaskHandler) EnqueueLLMTask(taskType string, payload map[string]interface{}) (*domainTask.Job, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.enqueueUseCase.Execute(ctx, domainTask.TaskType(taskType), payload)
}

// GetRecentTasks retorna la lista de tareas recientes para alimentar la UI
func (h *TaskHandler) GetRecentTasks(limit int) ([]domainTask.Job, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.repo.ListRecent(ctx, limit)
}

// GetTaskByID consulta el estado puntual de una tarea
func (h *TaskHandler) GetTaskByID(id string) (*domainTask.Job, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.repo.GetByID(ctx, id)
}
