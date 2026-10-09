package message

import (
	"context"
	"fmt"
	"time"

	"proyecto_ia/internal/domain"
)

type SendMessageStubUseCase struct {
	messageRepo domain.MessageRepository
	agentRepo   domain.AgentRepository
	metricRepo  domain.MetricRepository
}

func NewSendMessageStubUseCase(
	messageRepo domain.MessageRepository,
	agentRepo domain.AgentRepository,
	metricRepo domain.MetricRepository,
) *SendMessageStubUseCase {
	return &SendMessageStubUseCase{
		messageRepo: messageRepo,
		agentRepo:   agentRepo,
		metricRepo:  metricRepo,
	}
}

type SendMessageInput struct {
	SessionID string `json:"sessionId"`
	AgentID   string `json:"agentId"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
}

func (uc *SendMessageStubUseCase) Execute(ctx context.Context, input SendMessageInput) (*domain.Message, error) {
	// 1. Guardar mensaje del usuario
	userMsgID := fmt.Sprintf("msg-%d-u", time.Now().UnixNano())
	userSender := input.Sender
	if userSender == "" {
		userSender = "Developer"
	}
	userMsg, err := domain.NewMessage(userMsgID, input.SessionID, "", domain.RoleUser, userSender, input.Content, nil)
	if err != nil {
		return nil, err
	}
	if err := uc.messageRepo.Save(ctx, userMsg); err != nil {
		return nil, err
	}

	// 2. Obtener detalles del agente para la respuesta
	agentName := "Hub Orchestrator"
	agentProvider := "google"
	agentModel := "gemini-1.5-flash"
	agentRole := "IA Assistant"

	if input.AgentID != "" {
		if ag, err := uc.agentRepo.FindByID(ctx, input.AgentID); err == nil && ag != nil {
			agentName = ag.Name
			agentProvider = ag.Provider
			agentModel = ag.Model
			agentRole = ag.Role
		}
	}

	// 3. Crear telemetría de stub
	metricID := fmt.Sprintf("met-%d", time.Now().UnixNano())
	promptTokens := len(input.Content)/4 + 25
	completionTokens := 85
	totalLatency := int64(320)
	ttft := int64(95)
	finishReason := "stop"

	metric := domain.NewMetric(metricID, agentProvider, agentModel, totalLatency, ttft, promptTokens, completionTokens, finishReason)
	_ = uc.metricRepo.Save(ctx, metric)

	// 4. Crear respuesta del agente con stub inteligente
	agentReplyID := fmt.Sprintf("msg-%d-a", time.Now().UnixNano())
	replyText := fmt.Sprintf(
		"[%s - %s]\nProcesando tu solicitud en el Hub.\n\nContenido analizado: \"%s\"\n\nTodos los subsistemas se encuentran sincronizados según Clean Architecture. Listo para la ejecución de pipelines y conectores MCP.",
		agentName, agentRole, input.Content,
	)

	agentMsg, err := domain.NewMessage(agentReplyID, input.SessionID, input.AgentID, domain.RoleAssistant, agentName, replyText, metric)
	if err != nil {
		return nil, err
	}
	if err := uc.messageRepo.Save(ctx, agentMsg); err != nil {
		return nil, err
	}

	return agentMsg, nil
}

type ListMessagesUseCase struct {
	messageRepo domain.MessageRepository
}

func NewListMessagesUseCase(messageRepo domain.MessageRepository) *ListMessagesUseCase {
	return &ListMessagesUseCase{messageRepo: messageRepo}
}

func (uc *ListMessagesUseCase) Execute(ctx context.Context, sessionID string) ([]domain.Message, error) {
	return uc.messageRepo.ListBySession(ctx, sessionID)
}
