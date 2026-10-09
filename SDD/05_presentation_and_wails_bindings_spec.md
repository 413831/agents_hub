# SDD-05: Especificación de Presentación y Bindings de Wails (Presentation Spec)

La capa de **Interfaces** (`internal/interfaces/wails`) y la capa de **Presentación** (`frontend/`) definen los puntos de entrada para la interacción del usuario mediante la interfaz gráfica de escritorio.

---

## 1. Controladores de Entrada de Wails (`internal/interfaces/wails`)

Los controladores registrados en el `Bind: []interface{}` de `wails.Run` son expuestos automáticamente al entorno JavaScript del navegador a través de `window.go.wails.<Handler>`.

| Handler | Métodos Expuestos al Frontend | Responsabilidad |
| :--- | :--- | :--- |
| **`App`** | `Startup(ctx)`, `Shutdown(ctx)`, `Greet(name)` | Ciclo de vida de la ventana y puente de eventos `runtime.EventsEmit`. |
| **`LLMHandler`** | `GenerateCompletion(req)`, `StreamPrompt(sessionID, req)` | Entrada IPC para inferencia de modelos y canalización de streaming. |
| **`TaskHandler`** | `EnqueueLLMTask(taskType, payload)`, `GetRecentTasks(limit)`, `GetTaskByID(id)` | Entrada IPC para encolamiento asíncrono y consulta de jobs. |
| **`GitHandler`** | `InspectRepository(repoPath)`, `CommitChanges(repoPath, message, author, email)` | Entrada IPC para análisis de repositorio y operaciones de commit. |
| **`MetricsHandler`** | `GetMetricsSnapshot()` | Consulta del snapshot consolidado de métricas para renderizado en UI. |

---

## 2. Puente de Eventos Reactivos (Event Bridge)

Para evitar que el frontend deba realizar peticiones periódicas (polling) sobre tareas de larga duración:
1. El backend invoca `app.Emit(eventName, data)`.
2. Wails transmite el payload mediante el canal de eventos bidireccional nativo (`runtime.EventsOn`).
3. El frontend actualiza reactivamente el estado de la UI al recibir eventos como `task:updated`.

---

## 3. Vistas del Frontend (`frontend/src/`)

El frontend está estructurado en 4 vistas independientes sin lógica de negocio acoplada:

1. **`LLM Studio` (`#view-chat`):**
   - Selector de proveedor (`OpenAI`, `Gemini`).
   - Selector de modelo (`gpt-4o`, `gemini-1.5-pro`, etc.).
   - Panel conversacional con burbujas de chat estilizadas.
   - Caja de entrada y botón de envío de prompts.

2. **`Message Broker Tasks` (`#view-tasks`):**
   - Botón para encolar nuevas tareas pesadas en segundo plano.
   - Tabla con `Job ID`, `Tipo de Tarea`, `Estado` (`pending`, `running`, `completed`), `Progreso %` y `Resultado`.

3. **`Git Repositories` (`#view-git`):**
   - Campo para ingresar la ruta de un repositorio local.
   - Tabla de estado de archivos modificados en el árbol de trabajo (`Working Tree`).
   - Tabla de historial de commits recientes con hash, autor y mensaje.

4. **`Metrics & Grafana` (`#view-metrics`):**
   - 4 tarjetas KPI en tiempo real: Peticiones Totales, Tokens Procesados, Latencia Promedio y Workers Activos.
   - Especificación y visualizador del endpoint HTTP local de Prometheus (`http://localhost:2112/metrics`).
