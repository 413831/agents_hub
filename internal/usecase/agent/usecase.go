package agent

import (
	"context"

	"proyecto_ia/internal/domain"
)

type ListAgentsUseCase struct {
	repo domain.AgentRepository
}

func NewListAgentsUseCase(repo domain.AgentRepository) *ListAgentsUseCase {
	return &ListAgentsUseCase{repo: repo}
}

func (uc *ListAgentsUseCase) Execute(ctx context.Context) ([]domain.Agent, error) {
	return uc.repo.List(ctx)
}
