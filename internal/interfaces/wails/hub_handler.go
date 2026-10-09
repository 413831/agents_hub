package wails

import (
	"context"

	"proyecto_ia/internal/domain"
	agentUC "proyecto_ia/internal/usecase/agent"
	connectorUC "proyecto_ia/internal/usecase/connector"
	messageUC "proyecto_ia/internal/usecase/message"
	projectUC "proyecto_ia/internal/usecase/project"
	sessionUC "proyecto_ia/internal/usecase/session"
)

// HubHandler expone los casos de uso principales del Hub a través del IPC de Wails
type HubHandler struct {
	listProjectsUC    *projectUC.ListProjectsUseCase
	createProjectUC   *projectUC.CreateProjectUseCase
	listSessionsUC    *sessionUC.ListSessionsUseCase
	createSessionUC   *sessionUC.CreateSessionUseCase
	listMessagesUC    *messageUC.ListMessagesUseCase
	sendMessageStubUC *messageUC.SendMessageStubUseCase
	listConnectorsUC  *connectorUC.ListConnectorsUseCase
	toggleConnectorUC *connectorUC.ToggleConnectorUseCase
	listAgentsUC      *agentUC.ListAgentsUseCase
	metricRepo        domain.MetricRepository
	app               *App
}

func NewHubHandler(
	listProjectsUC *projectUC.ListProjectsUseCase,
	createProjectUC *projectUC.CreateProjectUseCase,
	listSessionsUC *sessionUC.ListSessionsUseCase,
	createSessionUC *sessionUC.CreateSessionUseCase,
	listMessagesUC *messageUC.ListMessagesUseCase,
	sendMessageStubUC *messageUC.SendMessageStubUseCase,
	listConnectorsUC *connectorUC.ListConnectorsUseCase,
	toggleConnectorUC *connectorUC.ToggleConnectorUseCase,
	listAgentsUC *agentUC.ListAgentsUseCase,
	metricRepo domain.MetricRepository,
	app *App,
) *HubHandler {
	return &HubHandler{
		listProjectsUC:    listProjectsUC,
		createProjectUC:   createProjectUC,
		listSessionsUC:    listSessionsUC,
		createSessionUC:   createSessionUC,
		listMessagesUC:    listMessagesUC,
		sendMessageStubUC: sendMessageStubUC,
		listConnectorsUC:  listConnectorsUC,
		toggleConnectorUC: toggleConnectorUC,
		listAgentsUC:      listAgentsUC,
		metricRepo:        metricRepo,
		app:               app,
	}
}

// ListProjects retorna todos los proyectos registrados
func (h *HubHandler) ListProjects() ([]domain.Project, error) {
	ctx := context.Background()
	return h.listProjectsUC.Execute(ctx)
}

// CreateProject registra un nuevo proyecto de trabajo
func (h *HubHandler) CreateProject(name, description string) (*domain.Project, error) {
	ctx := context.Background()
	proj, err := h.createProjectUC.Execute(ctx, projectUC.CreateProjectInput{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	h.app.Emit("project:created", proj)
	return proj, nil
}

// ListSessions retorna las sesiones correspondientes a un proyecto
func (h *HubHandler) ListSessions(projectID string) ([]domain.Session, error) {
	ctx := context.Background()
	return h.listSessionsUC.Execute(ctx, projectID)
}

// CreateSession crea una nueva sesión de chat interactiva
func (h *HubHandler) CreateSession(projectID, title string, agentIDs []string) (*domain.Session, error) {
	ctx := context.Background()
	sess, err := h.createSessionUC.Execute(ctx, sessionUC.CreateSessionInput{
		ProjectID: projectID,
		Title:     title,
		AgentIDs:  agentIDs,
	})
	if err != nil {
		return nil, err
	}
	h.app.Emit("session:created", sess)
	return sess, nil
}

// GetSessionMessages obtiene el historial de mensajes de una sesión
func (h *HubHandler) GetSessionMessages(sessionID string) ([]domain.Message, error) {
	ctx := context.Background()
	return h.listMessagesUC.Execute(ctx, sessionID)
}

// SendMessage envía un mensaje de usuario y genera la respuesta con stub y telemetría
func (h *HubHandler) SendMessage(sessionID, agentID, sender, content string) (*domain.Message, error) {
	ctx := context.Background()
	agentMsg, err := h.sendMessageStubUC.Execute(ctx, messageUC.SendMessageInput{
		SessionID: sessionID,
		AgentID:   agentID,
		Sender:    sender,
		Content:   content,
	})
	if err != nil {
		return nil, err
	}
	h.app.Emit("message:received", agentMsg)
	return agentMsg, nil
}

// ListConnectors retorna la lista de servidores MCP, skills y plugins
func (h *HubHandler) ListConnectors() ([]domain.Connector, error) {
	ctx := context.Background()
	return h.listConnectorsUC.Execute(ctx)
}

// ToggleConnector activa o desactiva un conector MCP/Skill
func (h *HubHandler) ToggleConnector(connectorID string) (*domain.Connector, error) {
	ctx := context.Background()
	conn, err := h.toggleConnectorUC.Execute(ctx, connectorID)
	if err != nil {
		return nil, err
	}
	h.app.Emit("connector:updated", conn)
	return conn, nil
}

// ListAgents retorna todos los agentes disponibles en el hub
func (h *HubHandler) ListAgents() ([]domain.Agent, error) {
	ctx := context.Background()
	return h.listAgentsUC.Execute(ctx)
}

// GetMetrics retorna las métricas de telemetría recientes
func (h *HubHandler) GetMetrics(limit int) ([]domain.Metric, error) {
	ctx := context.Background()
	if limit <= 0 {
		limit = 50
	}
	return h.metricRepo.ListRecent(ctx, limit)
}
