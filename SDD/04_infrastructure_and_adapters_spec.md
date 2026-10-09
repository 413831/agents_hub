# SDD-04: Especificación de Infraestructura y Adaptadores (Infrastructure Spec)

La capa de **Infraestructura** (`internal/infrastructure`) contiene los adaptadores concretos que implementan los puertos definidos por el Dominio y los Casos de Uso.

---

## 1. Adaptadores LLM (`internal/infrastructure/llm`)

### 1.1 `ProviderFactory` (Patrón Factory & Registry)
- Implementa `ProviderResolver`.
- Mantiene un registro concurrente seguro (`sync.RWMutex`) de instancias `LLMProvider`.
- Permite registrar o intercambiar proveedores dinámicamente mediante `Register(provider)`.

### 1.2 `OpenAIAdapter` (Patrón Strategy/Adapter)
- Encapsula la comunicación HTTP con `https://api.openai.com/v1/chat/completions`.
- Mapea códigos de estado HTTP a errores canónicos de dominio:
  - `401 Unauthorized` $\to$ `ErrInvalidAPIKey`
  - `429 Too Many Requests` $\to$ `ErrRateLimitExceeded`
- Soporta modo de simulación amigable si no se define `OPENAI_API_KEY`.

### 1.3 `GeminiAdapter` (Patrón Strategy/Adapter)
- Encapsula la comunicación con la API de Google Generative AI (`https://generativelanguage.googleapis.com/v1beta`).
- Transforma los mensajes canónicos (`system`, `user`, `assistant`) al esquema de contenidos y partes de Gemini (`user`, `model`).
- Soporta simulación si no se define `GEMINI_API_KEY`.

### 1.4 `TelemetryDecorator` (Patrón Decorator)
- Envuelve cualquier instancia de `LLMProvider`.
- Mide la duración exacta de la ejecución con `time.Since(start)` y extrae la cantidad de tokens consumidos.
- Registra las métricas en `MetricsRecorder` antes de retornar el resultado al caso de uso, **sin modificar el código del proveedor**.

---

## 2. Message Broker & Background Workers (`internal/infrastructure/broker` & `worker`)

### 2.1 `InProcessBroker` (Patrón Publisher/Subscriber)
- Broker embebido de alto rendimiento basado en canales concurrentes de Go (`chan struct`) y goroutines de despacho.
- **Ventaja de diseño:** Permite ejecución monolítica autónoma con cero instalación para el usuario final (no requiere servicio de Redis ni servidor externo).
- Implementa `domainTask.TaskBroker`.

### 2.2 `WorkerPool` (Pool Concurrente en Segundo Plano)
- Inicializa un conjunto parametrizable de goroutines trabajadoras (por defecto 3).
- Escucha tópicos como `tasks.task:llm:inference` y `tasks.task:git:clone`.
- Emite eventos reactivos a la UI (`task:updated`) ante transiciones de estado (`pending` $\to$ `running` $\to$ `completed` / `failed`) con actualización de porcentajes de progreso.

### 2.3 `InMemoryJobRepository`
- Almacenamiento concurrente protegido con `sync.RWMutex` para persistencia en memoria y consulta rápida de estados por ID o listado histórico ordenado cronológicamente.

---

## 3. Observabilidad con Prometheus & Grafana (`internal/infrastructure/telemetry`)

### 3.1 `PrometheusMetricsManager`
- Contadores atómicos (`sync/atomic`) y métricas protegidas:
  - `llm_requests_total`: Total de peticiones por proveedor y estado.
  - `llm_tokens_consumed_total`: Tokens totales procesados.
  - `llm_request_latency_avg_ms`: Latencia promedio acumulada.
  - `background_active_workers`: Cantidad de workers concurrentes activos.
- **Servidor HTTP Embebido:** Levanta un listener HTTP en `:2112/metrics` sirviendo métricas en formato estándar de Prometheus (versión 0.0.4) para scraping desde Prometheus o Grafana Agent / Alloy.

---

## 4. Adaptador Git (`internal/infrastructure/git`)

### 4.1 `NativeGitAdapter`
- Implementa `domainGit.GitClient`.
- Ejecuta operaciones de control de versiones locales (`git status --porcelain`, `git log`, `git commit`, `git clone`, `git branch`).
- Realiza fallbacks seguros cuando se abre un directorio que aún no está versionado.
