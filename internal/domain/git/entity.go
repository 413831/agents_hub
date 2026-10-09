package git

import "time"

// Repository representa un repositorio de código gestionado
type Repository struct {
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	RemoteURL     string    `json:"remote_url"`
	CurrentBranch string    `json:"current_branch"`
	IsClean       bool      `json:"is_clean"`
	LastSync      time.Time `json:"last_sync"`
}

// FileStatus estado de un archivo dentro del árbol de trabajo
type FileStatus struct {
	Path     string `json:"path"`
	Staging  string `json:"staging"` // 'M', 'A', 'D', '?'
	Worktree string `json:"worktree"`
}

// CommitInfo información de un commit
type CommitInfo struct {
	Hash        string    `json:"hash"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Message     string    `json:"message"`
	Date        time.Time `json:"date"`
}
