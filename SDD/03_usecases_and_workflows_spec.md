# SDD-03: Especificación de Casos de Uso y Flujos (Use Cases & Workflows)

La capa de **Casos de Uso** (`internal/usecase`) orquesta el flujo de datos entre las entidades de dominio y los adaptadores de infraestructura para cumplir los requerimientos de la aplicación.

---

## 1. Casos de Uso Implementados

### 1.1 `GenerateCompletionUseCase`
- **Ubicación:** `internal/usecase/llm/generate_completion.go`
- **Responsabilidad:** Valida la petición de inferencia, resuelve el proveedor adecuado a través de la interfaz `ProviderResolver`, ejecuta la llamada calculando la latencia en milisegundos (`DurationMs`) y retorna la respuesta canónica.
- **Flujo:**
  ```mermaid
  sequenceDiagram
      participant Handler as Wails LLMHandler
      participant UC as GenerateCompletionUseCase
      participant Resolver as ProviderResolver (Factory)
      participant Provider as LLMProvider (Adapter)

      Handler->>UC: Execute(ctx, req)
      UC->>Resolver: Get(req.Provider)
      Resolver-->>UC: instance (OpenAI/Gemini)
      UC->>Provider: GenerateCompletion(ctx, req)
      Provider-->>UC: *CompletionResponse
      UC-->>Handler: *CompletionResponse (con DurationMs)
  ```

### 1.2 `StreamCompletionUseCase`
- **Ubicación:** `internal/usecase/llm/stream_completion.go`
- **Responsabilidad:** Canaliza el streaming de tokens mediante el callback `onChunk`, propagando la cancelación de `context.Context` si el usuario cancela la sesión.

### 1.3 `EnqueueTaskUseCase`
- **Ubicación:** `internal/usecase/task/enqueue_task.go`
- **Responsabilidad:** Crea una entidad `Job` con ID único, estado inicial `pending`, la persiste en el repositorio de trabajos y la publica en el topic correspondiente (`tasks.<tipo>`) en el Message Broker.
- **Flujo de Tarea en Segundo Plano:**
  ```mermaid
  sequenceDiagram
      actor User as Usuario
      participant UI as Frontend
      participant Handler as TaskHandler
      participant UC as EnqueueTaskUseCase
      participant Repo as JobRepository
      participant Broker as TaskBroker
      participant Pool as WorkerPool
      participant Emitter as UIEventEmitter

      User->>UI: Clic "Encolar Tarea"
      UI->>Handler: EnqueueLLMTask(tipo, payload)
      Handler->>UC: Execute(ctx, tipo, payload)
      UC->>Repo: Save(job [pending])
      UC->>Broker: Publish("tasks.llm", job)
      UC-->>Handler: *Job
      Handler-->>UI: Retorna Job inmediatamente (No bloquea)

      Broker->>Pool: Despacha job a worker libre
      Pool->>Repo: UpdateStatus(running, 10%)
      Pool->>Emitter: Emit("task:updated", job)
      Emitter-->>UI: Actualiza fila en tabla

      Pool->>Pool: Ejecuta trabajo pesado...
      Pool->>Repo: UpdateStatus(completed, 100%)
      Pool->>Emitter: Emit("task:updated", job)
      Emitter-->>UI: Actualiza fila a "completed"
  ```

### 1.4 `InspectRepoUseCase`
- **Ubicación:** `internal/usecase/git/inspect_repo.go`
- **Responsabilidad:** Abre una ruta local, recupera los archivos modificados (`GetStatus`) y los últimos commits (`GetHistory`), consolidándolos en una estructura de presentación `RepoDetails`.
