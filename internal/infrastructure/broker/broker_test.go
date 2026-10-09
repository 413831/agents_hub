package broker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	domainTask "proyecto_ia/internal/domain/task"
	infraBroker "proyecto_ia/internal/infrastructure/broker"
)

func TestInProcessBroker_PublishSubscribe(t *testing.T) {
	broker := infraBroker.NewInProcessBroker(10)
	defer broker.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	topic := "tasks.test"
	var receivedJob domainTask.Job

	err := broker.Subscribe(ctx, topic, func(job domainTask.Job) {
		receivedJob = job
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}

	testJob := domainTask.Job{
		ID:     "job-999",
		Type:   domainTask.TypeLLMInference,
		Status: domainTask.StatusPending,
	}

	err = broker.Publish(ctx, topic, testJob)
	if err != nil {
		t.Fatalf("Publish error: %v", err)
	}

	wg.Wait()

	if receivedJob.ID != "job-999" {
		t.Errorf("Expected job ID job-999, got %s", receivedJob.ID)
	}
}
