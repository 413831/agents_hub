package git

import "context"

// GitClient puerto de dominio para interactuar con sistemas de control de versiones
type GitClient interface {
	Clone(ctx context.Context, repoURL, targetDir string) (*Repository, error)
	Open(ctx context.Context, path string) (*Repository, error)
	GetStatus(ctx context.Context, repoPath string) ([]FileStatus, error)
	Commit(ctx context.Context, repoPath, message, authorName, authorEmail string) (*CommitInfo, error)
	GetBranches(ctx context.Context, repoPath string) ([]string, error)
	GetHistory(ctx context.Context, repoPath string, limit int) ([]CommitInfo, error)
}
