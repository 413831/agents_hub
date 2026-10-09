package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	domainLLM "proyecto_ia/internal/domain/llm"
	domainTask "proyecto_ia/internal/domain/task"
)

// UIEventEmitter interfaz para emitir eventos a la ventana de Wails
type UIEventEmitter interface {
	Emit(eventName string, data interface{})
}

// WorkerPool procesa tareas asíncronas en segundo plano
type WorkerPool struct {
	broker    domainTask.TaskBroker
	repo      domainTask.JobRepository
	emitter   UIEventEmitter
	llmResolv interface {
		Get(pType domainLLM.ProviderType) (domainLLM.LLMProvider, error)
	}
	workerCount int
	jobChan     chan domainTask.Job
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWorkerPool(
	broker domainTask.TaskBroker,
	repo domainTask.JobRepository,
	emitter UIEventEmitter,
	llmResolv interface {
		Get(pType domainLLM.ProviderType) (domainLLM.LLMProvider, error)
	},
	workerCount int,
) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		broker:      broker,
		repo:        repo,
		emitter:     emitter,
		llmResolv:   llmResolv,
		workerCount: workerCount,
		jobChan:     make(chan domainTask.Job, 100),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (wp *WorkerPool) Start() {
	// Suscribirse a tópicos de tareas
	topics := []string{
		fmt.Sprintf("tasks.%s", domainTask.TypeLLMInference),
		fmt.Sprintf("tasks.%s", domainTask.TypeGitClone),
		fmt.Sprintf("tasks.%s", domainTask.TypeGitAnalysis),
	}

	for _, t := range topics {
		_ = wp.broker.Subscribe(wp.ctx, t, func(job domainTask.Job) {
			select {
			case wp.jobChan <- job:
			case <-wp.ctx.Done():
			}
		})
	}

	// Iniciar goroutines del pool
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.workerRoutine(i)
	}
}

func (wp *WorkerPool) workerRoutine(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobChan:
			if !ok {
				return
			}
			wp.processJob(job)
		}
	}
}

func (wp *WorkerPool) processJob(job domainTask.Job) {
	// 1. Actualizar a RUNNING
	_ = wp.repo.UpdateStatus(wp.ctx, job.ID, domainTask.StatusRunning, 10, nil, "")
	job.Status = domainTask.StatusRunning
	job.Progress = 10
	if wp.emitter != nil {
		wp.emitter.Emit("task:updated", job)
	}

	var resultBytes []byte
	var execErr error

	// 2. Procesar según el tipo de tarea
	switch job.Type {
	case domainTask.TypeLLMInference:
		var req domainLLM.CompletionRequest
		if err := json.Unmarshal(job.Payload, &req); err != nil {
			execErr = fmt.Errorf("invalid llm payload: %w", err)
		} else {
			provider, err := wp.llmResolv.Get(req.Provider)
			if err != nil {
				execErr = err
			} else {
				// Simulación de pasos de progreso en background
				for p := 25; p <= 75; p += 25 {
					time.Sleep(300 * time.Millisecond)
					_ = wp.repo.UpdateStatus(wp.ctx, job.ID, domainTask.StatusRunning, p, nil, "")
					job.Progress = p
					if wp.emitter != nil {
						wp.emitter.Emit("task:updated", job)
					}
				}

				resp, err := provider.GenerateCompletion(wp.ctx, req)
				if err != nil {
					execErr = err
				} else {
					resultBytes, _ = json.Marshal(resp)
				}
			}
		}

	default:
		// Tareas genéricas simuladas
		time.Sleep(800 * time.Millisecond)
		resultBytes = []byte(fmt.Sprintf(`{"status":"success","processed_at":"%s"}`, time.Now().Format(time.RFC3339)))
	}

	// 3. Finalizar estado
	if execErr != nil {
		_ = wp.repo.UpdateStatus(wp.ctx, job.ID, domainTask.StatusFailed, job.Progress, nil, execErr.Error())
		job.Status = domainTask.StatusFailed
		job.Error = execErr.Error()
	} else {
		_ = wp.repo.UpdateStatus(wp.ctx, job.ID, domainTask.StatusCompleted, 100, resultBytes, "")
		job.Status = domainTask.StatusCompleted
		job.Progress = 100
		job.Result = resultBytes
	}

	if wp.emitter != nil {
		wp.emitter.Emit("task:updated", job)
	}
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobChan)
	wp.wg.Wait()
}
