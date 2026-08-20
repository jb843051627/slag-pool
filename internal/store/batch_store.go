package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type BatchStore struct {
	db    *sql.DB
	mu    sync.RWMutex
	cache map[int64][]model.CoolingBatch
}

func (s *BatchStore) GetByID(ctx context.Context, id int64) (*model.CoolingBatch, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, pool_id, batch_number, start_time, end_time, target_temp, actual_temp, volume, source, status, created_at FROM cooling_batches WHERE id = ?", id)
	var b model.CoolingBatch
	var startTime, endTime, createdAt string
	if err := row.Scan(&b.ID, &b.PoolID, &b.BatchNumber, &startTime, &endTime, &b.TargetTemp, &b.ActualTemp, &b.Volume, &b.Source, &b.Status, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan batch: %w", err)
	}
	b.StartTime, _ = time.Parse(time.RFC3339, startTime)
	if endTime != "" {
		et, _ := time.Parse(time.RFC3339, endTime)
		b.EndTime = &et
	}
	b.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &b, nil
}

func (s *BatchStore) ListByPool(ctx context.Context, poolID int64) ([]model.CoolingBatch, error) {
	s.mu.RLock()
	if cached, ok := s.cache[poolID]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, batch_number, start_time, end_time, target_temp, actual_temp, volume, source, status, created_at FROM cooling_batches WHERE pool_id = ?", poolID)
	if err != nil {
		return nil, fmt.Errorf("list batches by pool: %w", err)
	}
	defer rows.Close()
	batches, err := scanBatches(rows)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[poolID] = batches
	s.mu.Unlock()
	return batches, nil
}

func (s *BatchStore) Create(ctx context.Context, batch *model.CoolingBatch) (int64, error) {
	if batch == nil || batch.PoolID == 0 {
		return 0, ErrInvalidInput
	}
	var endTimeVal any
	if batch.EndTime != nil {
		endTimeVal = batch.EndTime.Format(time.RFC3339)
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO cooling_batches (pool_id, batch_number, start_time, end_time, target_temp, actual_temp, volume, source, status, created_at) VALUES (?,?,?,?,?,?,?,?,?,?)",
		batch.PoolID, batch.BatchNumber, batch.StartTime.Format(time.RFC3339), endTimeVal,
		batch.TargetTemp, batch.ActualTemp, batch.Volume, batch.Source, batch.Status, batch.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create batch: %w", err)
	}
	id, _ := res.LastInsertId()
	batch.ID = id
	s.mu.Lock()
	delete(s.cache, batch.PoolID)
	s.mu.Unlock()
	return id, nil
}

func (s *BatchStore) Update(ctx context.Context, batch *model.CoolingBatch) error {
	if batch == nil || batch.ID == 0 {
		return ErrInvalidInput
	}
	var endTimeVal any
	if batch.EndTime != nil {
		endTimeVal = batch.EndTime.Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx,
		"UPDATE cooling_batches SET pool_id=?, batch_number=?, start_time=?, end_time=?, target_temp=?, actual_temp=?, volume=?, source=?, status=? WHERE id=?",
		batch.PoolID, batch.BatchNumber, batch.StartTime.Format(time.RFC3339), endTimeVal,
		batch.TargetTemp, batch.ActualTemp, batch.Volume, batch.Source, batch.Status, batch.ID)
	if err != nil {
		return fmt.Errorf("update batch: %w", err)
	}
	s.mu.Lock()
	delete(s.cache, batch.PoolID)
	s.mu.Unlock()
	return nil
}

func (s *BatchStore) ListByDateRange(ctx context.Context, start, end time.Time) ([]model.CoolingBatch, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, batch_number, start_time, end_time, target_temp, actual_temp, volume, source, status, created_at FROM cooling_batches WHERE start_time >= ? AND start_time <= ?",
		start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("list batches by date range: %w", err)
	}
	defer rows.Close()
	return scanBatches(rows)
}

func scanBatches(rows *sql.Rows) ([]model.CoolingBatch, error) {
	var batches []model.CoolingBatch
	for rows.Next() {
		var b model.CoolingBatch
		var startTime, endTime, createdAt string
		if err := rows.Scan(&b.ID, &b.PoolID, &b.BatchNumber, &startTime, &endTime, &b.TargetTemp, &b.ActualTemp, &b.Volume, &b.Source, &b.Status, &createdAt); err != nil {
			return nil, fmt.Errorf("scan batch row: %w", err)
		}
		b.StartTime, _ = time.Parse(time.RFC3339, startTime)
		if endTime != "" {
			et, _ := time.Parse(time.RFC3339, endTime)
			b.EndTime = &et
		}
		b.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		batches = append(batches, b)
	}
	return batches, nil
}
