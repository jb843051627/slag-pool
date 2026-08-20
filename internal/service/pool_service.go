package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type PoolService struct {
	store *store.Store
}

func NewPoolService(s *store.Store) *PoolService {
	return &PoolService{store: s}
}

func (s *PoolService) Get(ctx context.Context, id int64) (*model.SlagPool, error) {
	p, err := s.store.Pools().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, store.ErrPoolNotFound
	}
	return p, nil
}

func (s *PoolService) List(ctx context.Context, limit, offset int) ([]model.SlagPool, error) {
	return s.store.Pools().List(ctx, limit, offset)
}

func (s *PoolService) Create(ctx context.Context, pool *model.SlagPool) (int64, error) {
	if pool == nil || pool.Name == "" {
		return 0, store.ErrInvalidInput
	}
	now := time.Now()
	if pool.CreatedAt.IsZero() {
		pool.CreatedAt = now
	}
	pool.UpdatedAt = now
	if pool.Status == "" {
		pool.Status = model.PoolStatusActive
	}
	return s.store.Pools().Create(ctx, pool)
}

func (s *PoolService) Update(ctx context.Context, pool *model.SlagPool) error {
	if pool == nil || pool.ID == 0 {
		return store.ErrInvalidInput
	}
	pool.UpdatedAt = time.Now()
	return s.store.Pools().Update(ctx, pool)
}

func (s *PoolService) Delete(ctx context.Context, id int64) error {
	return s.store.Pools().Delete(ctx, id)
}

func (s *PoolService) GetStatus(ctx context.Context, poolID int64) (*model.PoolStatus, error) {
	p, err := s.store.Pools().GetByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, store.ErrPoolNotFound
	}
	status := &model.PoolStatus{PoolID: poolID}
	readings, _ := s.store.Readings().GetLatest(ctx, poolID)
	for _, r := range readings {
		switch {
		case r.Unit == "C" || r.Unit == "°C":
			status.Temp = r.Value
		case r.Unit == "L/min" || r.Unit == "m³/h":
			status.Flow = r.Value
		case r.Unit == "kPa" || r.Unit == "bar":
			status.Pressure = r.Value
		}
	}
	alerts, _ := s.store.Alerts().ListByPool(ctx, poolID)
	for _, a := range alerts {
		if a.Status == model.AlertStatusActive {
			status.AlertCount++
		}
	}
	batches, _ := s.store.Batches().ListByPool(ctx, poolID)
	for _, b := range batches {
		if b.Status == model.BatchStatusCooling {
			status.ActiveBatch = b.BatchNumber
			break
		}
	}
	return status, nil
}

func (s *PoolService) Exists(ctx context.Context, id int64) error {
	p, err := s.store.Pools().GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("pool %d not found: %w", id, store.ErrPoolNotFound)
	}
	return nil
}
