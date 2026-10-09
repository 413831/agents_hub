package connector

import (
	"context"

	"proyecto_ia/internal/domain"
)

type ListConnectorsUseCase struct {
	repo domain.ConnectorRepository
}

func NewListConnectorsUseCase(repo domain.ConnectorRepository) *ListConnectorsUseCase {
	return &ListConnectorsUseCase{repo: repo}
}

func (uc *ListConnectorsUseCase) Execute(ctx context.Context) ([]domain.Connector, error) {
	return uc.repo.List(ctx)
}

type ToggleConnectorUseCase struct {
	repo domain.ConnectorRepository
}

func NewToggleConnectorUseCase(repo domain.ConnectorRepository) *ToggleConnectorUseCase {
	return &ToggleConnectorUseCase{repo: repo}
}

func (uc *ToggleConnectorUseCase) Execute(ctx context.Context, id string) (*domain.Connector, error) {
	conn, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var newStatus domain.ConnectorStatus
	if conn.Status == domain.ConnectorStatusConnected {
		newStatus = domain.ConnectorStatusDisabled
	} else {
		newStatus = domain.ConnectorStatusConnected
	}
	if err := uc.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}
	conn.Status = newStatus
	return conn, nil
}
