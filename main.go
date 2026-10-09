package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	infraBroker "proyecto_ia/internal/infrastructure/broker"
	infraGit "proyecto_ia/internal/infrastructure/git"
	infraLLM "proyecto_ia/internal/infrastructure/llm"
	infraTelemetry "proyecto_ia/internal/infrastructure/telemetry"
	infraWorker "proyecto_ia/internal/infrastructure/worker"
	interfacesWails "proyecto_ia/internal/interfaces/wails"
	usecaseGit "proyecto_ia/internal/usecase/git"
	usecaseLLM "proyecto_ia/internal/usecase/llm"
	usecaseTask "proyecto_ia/internal/usecase/task"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fmt.Println(">> Iniciando AI Studio (Clean Architecture + Wails Monolith)...")

	// ---------------------------------------------------------
	// 1. CAPA DE INFRAESTRUCTURA (Adapters & Drivers)
	// ---------------------------------------------------------
	// 1.1 Observabilidad & Métricas Prometheus (Puerto 2112 para scraping de Grafana)
	metricsMgr := infraTelemetry.NewPrometheusMetricsManager()
	metricsMgr.StartMetricsServer(2112)
	metricsMgr.SetActiveWorkers(3)

	// 1.2 Adaptadores LLM (OpenAI & Gemini) decorados con Telemetría
	openaiRaw := infraLLM.NewOpenAIAdapter(os.Getenv("OPENAI_API_KEY"))
	geminiRaw := infraLLM.NewGeminiAdapter(os.Getenv("GEMINI_API_KEY"))

	openaiDecorated := infraLLM.NewTelemetryDecorator(openaiRaw, metricsMgr)
	geminiDecorated := infraLLM.NewTelemetryDecorator(geminiRaw, metricsMgr)

	// Registry & Factory
	llmFactory := infraLLM.NewProviderFactory()
	llmFactory.Register(openaiDecorated)
	llmFactory.Register(geminiDecorated)

	// 1.3 Message Broker & Repositorio de Tareas en Background
	broker := infraBroker.NewInProcessBroker(100)
	defer broker.Close()
	jobRepo := infraBroker.NewInMemoryJobRepository()

	// 1.4 Adaptador Git Nativo
	gitClient := infraGit.NewNativeGitAdapter()

	// ---------------------------------------------------------
	// 2. CAPA DE INTERFACES: APP CONTROLLER Y EMISOR DE EVENTOS
	// ---------------------------------------------------------
	app := interfacesWails.NewApp()

	// 1.5 Worker Pool (Fase 3: Procesamiento asíncrono desacoplado)
	workerPool := infraWorker.NewWorkerPool(broker, jobRepo, app, llmFactory, 3)
	workerPool.Start()
	defer workerPool.Stop()

	// ---------------------------------------------------------
	// 3. CAPA DE CASOS DE USO (Application Core)
	// ---------------------------------------------------------
	genCompletionUC := usecaseLLM.NewGenerateCompletionUseCase(llmFactory)
	streamCompletionUC := usecaseLLM.NewStreamCompletionUseCase(llmFactory)
	enqueueTaskUC := usecaseTask.NewEnqueueTaskUseCase(broker, jobRepo)
	inspectRepoUC := usecaseGit.NewInspectRepoUseCase(gitClient)

	// ---------------------------------------------------------
	// 4. CAPA DE INTERFACES: CONTROLADORES WAILS (Inbound Adapters)
	// ---------------------------------------------------------
	llmHandler := interfacesWails.NewLLMHandler(genCompletionUC, streamCompletionUC, app)
	taskHandler := interfacesWails.NewTaskHandler(enqueueTaskUC, jobRepo, app)
	gitHandler := interfacesWails.NewGitHandler(inspectRepoUC, gitClient, app)
	metricsHandler := interfacesWails.NewMetricsHandler(metricsMgr)

	// ---------------------------------------------------------
	// 5. COMPOSITION ROOT: INICIALIZACIÓN DE WAILS RUNTIME
	// ---------------------------------------------------------
	err := wails.Run(&options.App{
		Title:  "AI Studio Monolith - Clean Architecture",
		Width:  1280,
		Height: 850,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1}, // Slate-900 elegante
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
			fmt.Println(">> Runtime de Wails activo. Event bus conectado.")
		},
		OnShutdown: func(ctx context.Context) {
			app.Shutdown(ctx)
			_ = metricsMgr.Stop(ctx)
		},
		Bind: []interface{}{
			app,
			llmHandler,
			taskHandler,
			gitHandler,
			metricsHandler,
		},
	})

	if err != nil {
		fmt.Printf("Error ejecutando la aplicación: %v\n", err)
	}
}
