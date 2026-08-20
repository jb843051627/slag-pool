package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type MaintenanceService struct {
	store *store.Store
	pools *PoolService
}

func NewMaintenanceService(s *store.Store, ps *PoolService) *MaintenanceService {
	return &MaintenanceService{store: s, pools: ps}
}

func (s *MaintenanceService) Create(ctx context.Context, task *model.MaintenanceTask) (int64, error) {
	if task == nil || task.PoolID == 0 {
		return 0, store.ErrInvalidInput
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.Status == "" {
		task.Status = model.MaintStatusPending
	}
	return s.store.Maintenance().Create(ctx, task)
}

func (s *MaintenanceService) Get(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	t, err := s.store.Maintenance().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, store.ErrMaintenanceNotFound
	}
	return t, nil
}

func (s *MaintenanceService) ListByPool(ctx context.Context, poolID int64) ([]model.MaintenanceTask, error) {
	return s.store.Maintenance().ListByPool(ctx, poolID)
}

func (s *MaintenanceService) Complete(ctx context.Context, id int64, cost float64) error {
	t, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	t.Status = model.MaintStatusCompleted
	t.CompletedDate = &now
	t.Cost = cost
	return s.store.Maintenance().Update(ctx, t)
}

func (s *MaintenanceService) Schedule(ctx context.Context, task *model.MaintenanceTask) (int64, error) {
	if task == nil || task.PoolID == 0 {
		return 0, store.ErrInvalidInput
	}
	task.Status = model.MaintStatusScheduled
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	return s.store.Maintenance().Create(ctx, task)
}
