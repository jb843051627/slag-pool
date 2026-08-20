package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type ReadingService struct {
	store    *store.Store
	sensors  *SensorService
}

func NewReadingService(s *store.Store, ss *SensorService) *ReadingService {
	return &ReadingService{store: s, sensors: ss}
}

func (s *ReadingService) Ingest(ctx context.Context, reading *model.SensorReading) (int64, error) {
	if reading == nil || reading.SensorID == 0 {
		return 0, store.ErrInvalidInput
	}
	if reading.Timestamp.IsZero() {
		reading.Timestamp = time.Now()
	}
	if reading.Quality == "" {
		reading.Quality = model.ReadingQualityGood
	}
	return s.store.Readings().Create(ctx, reading)
}

func (s *ReadingService) ListBySensor(ctx context.Context, sensorID int64, limit int) ([]model.SensorReading, error) {
	return s.store.Readings().ListBySensor(ctx, sensorID, int64(limit))
}

func (s *ReadingService) ListByPool(ctx context.Context, poolID int64, limit int) ([]model.SensorReading, error) {
	return s.store.Readings().ListByPool(ctx, poolID, int64(limit))
}

func (s *ReadingService) BatchIngest(ctx context.Context, batch *model.ReadingBatch) error {
	if batch == nil || len(batch.Readings) == 0 {
		return store.ErrInvalidInput
	}
	for i := range batch.Readings {
		if batch.Readings[i].PoolID == 0 {
			batch.Readings[i].PoolID = batch.PoolID
		}
		if batch.Readings[i].Timestamp.IsZero() {
			batch.Readings[i].Timestamp = time.Now()
		}
		if batch.Readings[i].Quality == "" {
			batch.Readings[i].Quality = model.ReadingQualityGood
		}
	}
	return s.store.Readings().BatchCreate(ctx, batch.Readings)
}

func (s *ReadingService) GetLatest(ctx context.Context, poolID int64) ([]model.SensorReading, error) {
	return s.store.Readings().GetLatest(ctx, poolID)
}
