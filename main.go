package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	infraBroker "proyecto_ia/internal/infrastructure/broker"
	infraGit "proyecto_ia/internal/infrastructure/git"
	infraLLM "proyecto_ia/internal/infrastructure/llm"
	infraMemory "proyecto_ia/internal/infrastructure/memory"
	infraTelemetry "proyecto_ia/internal/infrastructure/telemetry"
	infraWorker "proyecto_ia/internal/infrastructure/worker"
	interfacesWails "proyecto_ia/internal/interfaces/wails"
	usecaseAgent "proyecto_ia/internal/usecase/agent"
	usecaseConnector "proyecto_ia/internal/usecase/connector"
	usecaseGit "proyecto_ia/internal/usecase/git"
	usecaseLLM "proyecto_ia/internal/usecase/llm"
	usecaseMessage "proyecto_ia/internal/usecase/message"
	usecaseProject "proyecto_ia/internal/usecase/project"
	usecaseSession "proyecto_ia/internal/usecase/session"
	usecaseTask "proyecto_ia/internal/usecase/task"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fmt.Println(">> Iniciando Agent & LLM Hub (Clean Architecture + DDD + Wails)...")

	// ---------------------------------------------------------
	// 1. CAPA DE INFRAESTRUCTURA (Adapters, Repositories & Drivers)
	// ---------------------------------------------------------
	// 1.1 Observabilidad & Métricas Prometheus (Puerto 2112)
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

	// 1.5 Repositorios en Memoria de Dominio (Hub DDD)
	projectRepo := infraMemory.NewInMemoryProjectRepository()
	sessionRepo := infraMemory.NewInMemorySessionRepository()
	messageRepo := infraMemory.NewInMemoryMessageRepository()
	agentRepo := infraMemory.NewInMemoryAgentRepository()
	connectorRepo := infraMemory.NewInMemoryConnectorRepository()
	metricRepo := infraMemory.NewInMemoryMetricRepository()

	// ---------------------------------------------------------
	// 2. CAPA DE INTERFACES: APP CONTROLLER Y EMISOR DE EVENTOS
	// ---------------------------------------------------------
	app := interfacesWails.NewApp()

	// Worker Pool para procesamiento asíncrono
	workerPool := infraWorker.NewWorkerPool(broker, jobRepo, app, llmFactory, 3)
	workerPool.Start()
	defer workerPool.Stop()

	// ---------------------------------------------------------
	// 3. CAPA DE CASOS DE USO (Application Core)
	// ---------------------------------------------------------
	// Casos de uso originales
	genCompletionUC := usecaseLLM.NewGenerateCompletionUseCase(llmFactory)
	streamCompletionUC := usecaseLLM.NewStreamCompletionUseCase(llmFactory)
	enqueueTaskUC := usecaseTask.NewEnqueueTaskUseCase(broker, jobRepo)
	inspectRepoUC := usecaseGit.NewInspectRepoUseCase(gitClient)

	// Casos de uso del Hub (Fase 1 DDD)
	listProjectsUC := usecaseProject.NewListProjectsUseCase(projectRepo)
	createProjectUC := usecaseProject.NewCreateProjectUseCase(projectRepo)
	listSessionsUC := usecaseSession.NewListSessionsUseCase(sessionRepo)
	createSessionUC := usecaseSession.NewCreateSessionUseCase(sessionRepo, projectRepo)
	listMessagesUC := usecaseMessage.NewListMessagesUseCase(messageRepo)
	sendMessageStubUC := usecaseMessage.NewSendMessageStubUseCase(messageRepo, agentRepo, metricRepo)
	listConnectorsUC := usecaseConnector.NewListConnectorsUseCase(connectorRepo)
	toggleConnectorUC := usecaseConnector.NewToggleConnectorUseCase(connectorRepo)
	listAgentsUC := usecaseAgent.NewListAgentsUseCase(agentRepo)

	// ---------------------------------------------------------
	// 4. CAPA DE INTERFACES: CONTROLADORES WAILS (Inbound Adapters)
	// ---------------------------------------------------------
	llmHandler := interfacesWails.NewLLMHandler(genCompletionUC, streamCompletionUC, app)
	taskHandler := interfacesWails.NewTaskHandler(enqueueTaskUC, jobRepo, app)
	gitHandler := interfacesWails.NewGitHandler(inspectRepoUC, gitClient, app)
	metricsHandler := interfacesWails.NewMetricsHandler(metricsMgr)
	hubHandler := interfacesWails.NewHubHandler(
		listProjectsUC,
		createProjectUC,
		listSessionsUC,
		createSessionUC,
		listMessagesUC,
		sendMessageStubUC,
		listConnectorsUC,
		toggleConnectorUC,
		listAgentsUC,
		metricRepo,
		app,
	)

	// ---------------------------------------------------------
	// 5. COMPOSITION ROOT: INICIALIZACIÓN DE WAILS RUNTIME
	// ---------------------------------------------------------
	err := wails.Run(&options.App{
		Title:  "Agent & LLM Hub - Clean Architecture",
		Width:  1320,
		Height: 880,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1}, // Slate-900
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
			hubHandler,
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
