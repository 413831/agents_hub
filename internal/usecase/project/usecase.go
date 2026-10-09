package project

import (
	"context"
	"fmt"
	"time"

	"proyecto_ia/internal/domain"
)

type ListProjectsUseCase struct {
	repo domain.ProjectRepository
}

func NewListProjectsUseCase(repo domain.ProjectRepository) *ListProjectsUseCase {
	return &ListProjectsUseCase{repo: repo}
}

func (uc *ListProjectsUseCase) Execute(ctx context.Context) ([]domain.Project, error) {
	return uc.repo.List(ctx)
}

type CreateProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewCreateProjectUseCase(repo domain.ProjectRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo}
}

type CreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, input CreateProjectInput) (*domain.Project, error) {
	id := fmt.Sprintf("proj-%d", time.Now().UnixNano())
	proj, err := domain.NewProject(id, input.Name, input.Description)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, proj); err != nil {
		return nil, err
	}
	return proj, nil
}
