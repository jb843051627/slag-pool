package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/slag-pool/internal/model"
)

func TestBug05_BatchIngestHonorsContextCancel(t *testing.T) {
	st := newTestStore(t)
	svc := NewReadingService(st, NewSensorService(st, NewPoolService(st)))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	batch := &model.ReadingBatch{PoolID: 1, Readings: []model.SensorReading{{SensorID: 1, PoolID: 1, Value: 1.0, Unit: "c"}}}
	err := svc.BatchIngest(ctx, batch)
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("BatchIngest: expected context.Canceled, got %v", err)
	}
}