package wails

import (
	"context"

	domainGit "proyecto_ia/internal/domain/git"
	usecaseGit "proyecto_ia/internal/usecase/git"
)

// GitHandler expone operaciones con repositorios Git al frontend
type GitHandler struct {
	inspectUseCase *usecaseGit.InspectRepoUseCase
	gitClient      domainGit.GitClient
	app            *App
}

func NewGitHandler(
	inspectUseCase *usecaseGit.InspectRepoUseCase,
	gitClient domainGit.GitClient,
	app *App,
) *GitHandler {
	return &GitHandler{
		inspectUseCase: inspectUseCase,
		gitClient:      gitClient,
		app:            app,
	}
}

// InspectRepository analiza el estado actual, archivos modificados y commits recientes
func (h *GitHandler) InspectRepository(repoPath string) (*usecaseGit.RepoDetails, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.inspectUseCase.Execute(ctx, repoPath)
}

// CommitChanges realiza commit en el repositorio indicado
func (h *GitHandler) CommitChanges(repoPath, message, author, email string) (*domainGit.CommitInfo, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.gitClient.Commit(ctx, repoPath, message, author, email)
}
