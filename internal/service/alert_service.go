package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type AlertService struct {
	store *store.Store
}

func NewAlertService(s *store.Store) *AlertService {
	return &AlertService{store: s}
}

func (s *AlertService) Create(ctx context.Context, alert *model.Alert) (int64, error) {
	if alert == nil || alert.PoolID == 0 {
		return 0, store.ErrInvalidInput
	}
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now()
	}
	if alert.Status == "" {
		alert.Status = model.AlertStatusActive
	}
	return s.store.Alerts().Create(ctx, alert)
}

func (s *AlertService) Get(ctx context.Context, id int64) (*model.Alert, error) {
	a, err := s.store.Alerts().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, store.ErrAlertNotFound
	}
	return a, nil
}

func (s *AlertService) ListByPool(ctx context.Context, poolID int64) ([]model.Alert, error) {
	return s.store.Alerts().ListByPool(ctx, poolID)
}

func (s *AlertService) ListByStatus(ctx context.Context, status string) ([]model.Alert, error) {
	return s.store.Alerts().ListByStatus(ctx, status)
}

func (s *AlertService) Acknowledge(ctx context.Context, id int64, routedTo string) error {
	return s.store.Alerts().Acknowledge(ctx, id, routedTo)
}

func (s *AlertService) Resolve(ctx context.Context, id int64) error {
	a, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	a.Status = model.AlertStatusResolved
	a.ResolvedAt = &now
	return s.store.Alerts().Update(ctx, a)
}
