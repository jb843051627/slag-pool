package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type BatchService struct {
	store *store.Store
	pools *PoolService
}

func NewBatchService(s *store.Store, ps *PoolService) *BatchService {
	return &BatchService{store: s, pools: ps}
}

func (s *BatchService) Get(ctx context.Context, id int64) (*model.CoolingBatch, error) {
	b, err := s.store.Batches().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, store.ErrBatchNotFound
	}
	return b, nil
}

func (s *BatchService) ListByPool(ctx context.Context, poolID int64) ([]model.CoolingBatch, error) {
	return s.store.Batches().ListByPool(ctx, poolID)
}

func (s *BatchService) Create(ctx context.Context, batch *model.CoolingBatch) (int64, error) {
	if batch == nil || batch.PoolID == 0 {
		return 0, store.ErrInvalidInput
	}
	if batch.CreatedAt.IsZero() {
		batch.CreatedAt = time.Now()
	}
	if batch.Status == "" {
		batch.Status = model.BatchStatusPending
	}
	return s.store.Batches().Create(ctx, batch)
}

func (s *BatchService) Start(ctx context.Context, id int64) error {
	b, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	b.Status = model.BatchStatusCooling
	b.StartTime = time.Now()
	return s.store.Batches().Update(ctx, b)
}

func (s *BatchService) Complete(ctx context.Context, id int64, actualTemp float64) error {
	b, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	b.Status = model.BatchStatusComplete
	b.ActualTemp = actualTemp
	b.EndTime = &now
	return s.store.Batches().Update(ctx, b)
}
