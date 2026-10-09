# SDD-01: Especificación de Arquitectura del Sistema (Clean Architecture Monolith)

## 1. Visión General del Sistema
**AI Studio** es una aplicación monolítica de escritorio multiplataforma construida con **Wails v2** (Go + Frontend Web) que aplica rigurosamente **Clean Architecture** (Arquitectura Limpia / Hexagonal / Ports & Adapters).

El sistema desacopla completamente el núcleo de negocio de las tecnologías externas (modelos de lenguaje, message brokers, sistemas de archivos Git y mecanismos de observabilidad), permitiendo que cualquier componente de infraestructura sea reemplazado o probado de forma aislada sin tocar las reglas de dominio.

---

## 2. Principios y Regla de Dependencia

```
+-------------------------------------------------------------------------+
|                  CAPA 4: INFRAESTRUCTURA & FRAMEWORKS                   |
|   (OpenAI, Gemini SDK, go-git, NATS/Redis, Prometheus, Wails Runtime)   |
|   +-----------------------------------------------------------------+   |
|   |                 CAPA 3: INTERFACES & ADAPTADORES                |   |
|   |         (Wails IPC Handlers, View Models, Event Bridges)        |   |
|   |   +---------------------------------------------------------+   |   |
|   |   |               CAPA 2: CASOS DE USO (USE CASES)          |   |   |
|   |   |          (GenerateCompletion, EnqueueJob, InspectRepo)  |   |   |
|   |   |   +-------------------------------------------------+   |   |   |
|   |   |   |              CAPA 1: DOMINIO (CORE)             |   |   |   |
|   |   |   |      (Entidades, Objetos de Valor, Puertos)     |   |   |   |
|   |   |   +-------------------------------------------------+   |   |   |
|   |   +---------------------------------------------------------+   |   |
|   +-----------------------------------------------------------------+   |
+-------------------------------------------------------------------------+
```

### Reglas Inviolables
1. **Regla de Dependencia:** Las dependencias del código fuente solo apuntan hacia adentro. Las capas internas nunca importan paquetes de capas externas.
2. **Independencia Tecnológica:** El paquete `internal/domain` NO tiene dependencias de Wails, APIs REST, drivers SQL, NATS ni Prometheus. Solo utiliza la biblioteca estándar de Go (`context`, `time`, `errors`).
3. **Inversión de Dependencias (DIP):** Las capas de Casos de Uso y Dominio declaran interfaces abstractas (**Puertos**). Las implementaciones concretas (**Adaptadores**) se encuentran en `internal/infrastructure` e `internal/interfaces`.
4. **Composition Root:** Toda la inyección de dependencias se resuelve en un único punto de entrada: `main.go`.

---

## 3. Topología de Directorios Canónica

```
proyecto_ia/
├── SDD/                                # Especificaciones (Spec-Driven Development)
│   ├── 01_system_architecture_spec.md
│   ├── 02_domain_and_ports_spec.md
│   ├── 03_usecases_and_workflows_spec.md
│   ├── 04_infrastructure_and_adapters_spec.md
│   ├── 05_presentation_and_wails_bindings_spec.md
│   └── 06_roadmap_and_extension_guide.md
├── build/                              # Recursos de empaquetado Wails (iconos, binarios)
│   └── bin/
│       └── proyecto_ia.exe
├── frontend/                           # Presentación desacoplada (Vite + Vanilla JS / CSS)
│   ├── dist/                           # Bundle empaquetado para embed FS
│   ├── src/
│   │   ├── index.css                   # Sistema de diseño y temas oscuros
│   │   └── main.js                     # Controlador de vistas y puente IPC
│   ├── wailsjs/                        # Bindings autogenerados por Wails
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── internal/
│   ├── domain/                         # Entidades de negocio y Puertos (Core)
│   │   ├── llm/                        # Message, CompletionRequest, LLMProvider
│   │   ├── task/                       # Job, TaskStatus, TaskBroker, JobRepository
│   │   └── git/                        # Repository, CommitInfo, GitClient
│   ├── usecase/                        # Orquestadores de Casos de Uso
│   │   ├── llm/                        # GenerateCompletion, StreamCompletion
│   │   ├── task/                       # EnqueueTask
│   │   └── git/                        # InspectRepo
│   ├── interfaces/                     # Adaptadores de entrada (Wails IPC)
│   │   └── wails/                      # App, LLMHandler, TaskHandler, GitHandler, MetricsHandler
│   └── infrastructure/                 # Adaptadores de salida (Tecnología concreta)
│       ├── llm/                        # Factory, OpenAIAdapter, GeminiAdapter, TelemetryDecorator
│       ├── broker/                     # InProcessBroker, InMemoryJobRepository
│       ├── worker/                     # WorkerPool concurrente en background
│       ├── git/                        # NativeGitAdapter (o go-git/v5)
│       └── telemetry/                  # PrometheusMetricsManager & HTTP Server
├── main.go                             # Composition Root (Bootstrap de dependencias y Wails)
├── wails.json                          # Configuración del proyecto Wails
├── go.mod                              # Módulo Go
└── go.sum                              # Checksums de dependencias
```
