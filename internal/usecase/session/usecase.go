package session

import (
	"context"
	"fmt"
	"time"

	"proyecto_ia/internal/domain"
)

type CreateSessionUseCase struct {
	sessionRepo domain.SessionRepository
	projectRepo domain.ProjectRepository
}

func NewCreateSessionUseCase(sessionRepo domain.SessionRepository, projectRepo domain.ProjectRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
		projectRepo: projectRepo,
	}
}

type CreateSessionInput struct {
	ProjectID string   `json:"projectId"`
	Title     string   `json:"title"`
	AgentIDs  []string `json:"agentIds"`
}

func (uc *CreateSessionUseCase) Execute(ctx context.Context, input CreateSessionInput) (*domain.Session, error) {
	// Verificar que el proyecto existe
	if input.ProjectID != "" {
		if _, err := uc.projectRepo.FindByID(ctx, input.ProjectID); err != nil {
			return nil, err
		}
	}
	id := fmt.Sprintf("sess-%d", time.Now().UnixNano())
	sess, err := domain.NewSession(id, input.ProjectID, input.Title, input.AgentIDs)
	if err != nil {
		return nil, err
	}
	if err := uc.sessionRepo.Save(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

type ListSessionsUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewListSessionsUseCase(sessionRepo domain.SessionRepository) *ListSessionsUseCase {
	return &ListSessionsUseCase{sessionRepo: sessionRepo}
}

func (uc *ListSessionsUseCase) Execute(ctx context.Context, projectID string) ([]domain.Session, error) {
	return uc.sessionRepo.ListByProject(ctx, projectID)
}
