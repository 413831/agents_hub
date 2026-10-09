# SDD-02: Especificación de Dominio y Puertos (Domain & Ports Spec)

La capa de **Dominio** (`internal/domain`) encapsula las reglas de negocio de la empresa y define las interfaces abstractas (Puertos) que determinan cómo interactúa el sistema con el mundo exterior.

---

## 1. Módulo LLM (`internal/domain/llm`)

### 1.1 Entidades y Tipos
- **`ProviderType`**: `string` enum que define los proveedores soportados (`openai`, `gemini`, `mock`).
- **`Role`**: `string` enum (`system`, `user`, `assistant`).
- **`Message`**: Entidad básica de interacción en un diálogo:
  ```go
  type Message struct {
      Role      Role      `json:"role"`
      Content   string    `json:"content"`
      Timestamp time.Time `json:"timestamp"`
  }
  ```
- **`CompletionRequest`**: DTO canónico de petición de inferencia:
  ```go
  type CompletionRequest struct {
      Provider    ProviderType `json:"provider"`
      Model       string       `json:"model"`
      Messages    []Message    `json:"messages"`
      Temperature float32      `json:"temperature"`
      MaxTokens   int          `json:"max_tokens"`
  }
  ```
- **`CompletionResponse`**: DTO canónico de respuesta generada:
  ```go
  type CompletionResponse struct {
      Provider     ProviderType `json:"provider"`
      Model        string       `json:"model"`
      Content      string       `json:"content"`
      PromptTokens int          `json:"prompt_tokens"`
      CompTokens   int          `json:"completion_tokens"`
      TotalTokens  int          `json:"total_tokens"`
      DurationMs   int64        `json:"duration_ms"`
  }
  ```

### 1.2 Puertos (Interfaces Primarias)
- **`LLMProvider`**:
  ```go
  type ChunkCallback func(chunk string) error

  type LLMProvider interface {
      Type() ProviderType
      SupportedModels() []string
      GenerateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
      StreamCompletion(ctx context.Context, req CompletionRequest, onChunk ChunkCallback) error
  }
  ```

### 1.3 Contratos de Errores de Dominio
- `ErrProviderNotFound`: Se solicitó un motor no registrado.
- `ErrInvalidAPIKey`: Credencial inválida o no configurada.
- `ErrRateLimitExceeded`: Cuota o frecuencia agotada en el upstream.
- `ErrContextTooLong`: Límite de tokens de la ventana de contexto superado.
- `ErrProviderTimeout`: Tiempo de respuesta excedido.

---

## 2. Módulo de Tareas Asíncronas (`internal/domain/task`)

### 2.1 Entidades y Tipos
- **`TaskStatus`**: `pending` | `running` | `completed` | `failed`.
- **`TaskType`**: Identificador del tópico del trabajo (`task:llm:inference`, `task:git:clone`, `task:git:analysis`).
- **`Job`**:
  ```go
  type Job struct {
      ID        string          `json:"id"`
      Type      TaskType        `json:"type"`
      Status    TaskStatus      `json:"status"`
      Progress  int             `json:"progress"` // 0 a 100%
      Payload   json.RawMessage `json:"payload"`
      Result    json.RawMessage `json:"result,omitempty"`
      Error     string          `json:"error,omitempty"`
      CreatedAt time.Time       `json:"created_at"`
      UpdatedAt time.Time       `json:"updated_at"`
  }
  ```

### 2.2 Puertos
- **`TaskBroker`**:
  ```go
  type TaskBroker interface {
      Publish(ctx context.Context, topic string, job Job) error
      Subscribe(ctx context.Context, topic string, handler func(job Job)) error
      Close() error
  }
  ```
- **`JobRepository`**:
  ```go
  type JobRepository interface {
      Save(ctx context.Context, job Job) error
      GetByID(ctx context.Context, id string) (*Job, error)
      ListRecent(ctx context.Context, limit int) ([]Job, error)
      UpdateStatus(ctx context.Context, id string, status TaskStatus, progress int, result []byte, errStr string) error
  }
  ```

---

## 3. Módulo Git (`internal/domain/git`)

### 3.1 Entidades
- **`Repository`**: Metadatos de un repositorio (`Name`, `Path`, `RemoteURL`, `CurrentBranch`, `IsClean`, `LastSync`).
- **`FileStatus`**: Estado en el árbol de trabajo (`Path`, `Staging`, `Worktree`).
- **`CommitInfo`**: Registro de un commit (`Hash`, `AuthorName`, `AuthorEmail`, `Message`, `Date`).

### 3.2 Puertos
- **`GitClient`**:
  ```go
  type GitClient interface {
      Clone(ctx context.Context, repoURL, targetDir string) (*Repository, error)
      Open(ctx context.Context, path string) (*Repository, error)
      GetStatus(ctx context.Context, repoPath string) ([]FileStatus, error)
      Commit(ctx context.Context, repoPath, message, authorName, authorEmail string) (*CommitInfo, error)
      GetBranches(ctx context.Context, repoPath string) ([]string, error)
      GetHistory(ctx context.Context, repoPath string, limit int) ([]CommitInfo, error)
  }
  ```
