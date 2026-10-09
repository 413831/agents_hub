# SDD-06: Guía de Extensión y Próximos Pasos (Extension & Roadmap Guide)

Esta guía explica cómo continuar desarrollando sobre este monolito una vez copiado en tu repositorio de Git.

---

## 1. Cómo Agregar un Nuevo Proveedor de LLM (ej. Anthropic Claude u Ollama)

Gracias a Clean Architecture, no se necesita modificar ningún caso de uso existente ni el frontend:

1. **Crear el Adaptador en Infraestructura:**
   Crear un archivo en `internal/infrastructure/llm/anthropic_adapter.go` (o `ollama_adapter.go`) que implemente la interfaz `domainLLM.LLMProvider`:
   ```go
   type AnthropicAdapter struct { ... }
   func (a *AnthropicAdapter) Type() domainLLM.ProviderType { return "anthropic" }
   func (a *AnthropicAdapter) SupportedModels() []string { return []string{"claude-3-5-sonnet"} }
   func (a *AnthropicAdapter) GenerateCompletion(...) (*domainLLM.CompletionResponse, error) { ... }
   func (a *AnthropicAdapter) StreamCompletion(...) error { ... }
   ```

2. **Registrarlo en `main.go`:**
   ```go
   anthropicRaw := infraLLM.NewAnthropicAdapter(os.Getenv("ANTHROPIC_API_KEY"))
   anthropicDecorated := infraLLM.NewTelemetryDecorator(anthropicRaw, metricsMgr)
   llmFactory.Register(anthropicDecorated)
   ```

3. **Agregar la opción en el HTML:**
   En `frontend/index.html`, agregar `<option value="anthropic">Anthropic Claude</option>`. ¡Listo!

---

## 2. Cómo Conectar un Broker Externo (NATS / Redis)

El puerto `domainTask.TaskBroker` es agnóstico a la implementación. Para usar Redis:

1. Crear `internal/infrastructure/broker/redis_broker.go` utilizando `github.com/redis/go-redis/v9`.
2. Implementar `Publish`, `Subscribe` y `Close`.
3. En `main.go`, reemplazar `infraBroker.NewInProcessBroker(...)` por `infraBroker.NewRedisBroker("localhost:6379")`.

El resto del sistema seguirá funcionando idénticamente.

---

## 3. Persistencia con SQLite

Actualmente el repositorio de jobs opera en memoria (`InMemoryJobRepository`).
Para persistir en disco entre reinicios de la aplicación:
1. Crear `internal/infrastructure/storage/sqlite_job_repository.go` implementando `domainTask.JobRepository`.
2. Utilizar `modernc.org/sqlite` (Go puro, sin necesidad de CGO).
3. Inyectarlo en `main.go` en lugar de `NewInMemoryJobRepository()`.

---

## 4. Ejecución de Tests y Aseguramiento de Calidad

Para validar que cualquier cambio preserve los contratos de Clean Architecture:

```powershell
# Ejecutar todas las pruebas unitarias y de integración
go test -v ./...

# Validar compilación con tags de escritorio
go build -tags desktop -o app_test.exe .
```
