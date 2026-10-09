package git

import (
	"context"
	"fmt"

	domainGit "proyecto_ia/internal/domain/git"
)

// InspectRepoUseCase consulta el estado de un repositorio local
type InspectRepoUseCase struct {
	client domainGit.GitClient
}

func NewInspectRepoUseCase(client domainGit.GitClient) *InspectRepoUseCase {
	return &InspectRepoUseCase{
		client: client,
	}
}

type RepoDetails struct {
	Repo    *domainGit.Repository   `json:"repo"`
	Status  []domainGit.FileStatus  `json:"status"`
	History []domainGit.CommitInfo  `json:"history"`
}

func (uc *InspectRepoUseCase) Execute(ctx context.Context, repoPath string) (*RepoDetails, error) {
	repo, err := uc.client.Open(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", repoPath, err)
	}

	status, err := uc.client.GetStatus(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for %s: %w", repoPath, err)
	}

	history, err := uc.client.GetHistory(ctx, repoPath, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %w", repoPath, err)
	}

	return &RepoDetails{
		Repo:    repo,
		Status:  status,
		History: history,
	}, nil
}
