package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type QualityStore struct {
	db *sql.DB
}

func (s *QualityStore) Create(ctx context.Context, wq *model.WaterQuality) (int64, error) {
	if wq == nil || wq.PoolID == 0 {
		return 0, ErrInvalidInput
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO water_quality (pool_id, ph, turbidity, dissolved_oxygen, temperature, conductivity, suspended_solids, timestamp) VALUES (?,?,?,?,?,?,?,?)",
		wq.PoolID, wq.PH, wq.Turbidity, wq.DissolvedOxygen, wq.Temperature, wq.Conductivity, wq.SuspendedSolids, wq.Timestamp.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create water quality: %w", err)
	}
	id, _ := res.LastInsertId()
	wq.ID = id
	return id, nil
}

func (s *QualityStore) GetLatest(ctx context.Context, poolID int64) (*model.WaterQuality, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, pool_id, ph, turbidity, dissolved_oxygen, temperature, conductivity, suspended_solids, timestamp FROM water_quality WHERE pool_id = ? ORDER BY timestamp DESC LIMIT 1", poolID)
	var wq model.WaterQuality
	var ts string
	if err := row.Scan(&wq.ID, &wq.PoolID, &wq.PH, &wq.Turbidity, &wq.DissolvedOxygen, &wq.Temperature, &wq.Conductivity, &wq.SuspendedSolids, &ts); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan water quality: %w", err)
	}
	wq.Timestamp, _ = time.Parse(time.RFC3339, ts)
	return &wq, nil
}

func (s *QualityStore) ListByPool(ctx context.Context, poolID int64, limit int) ([]model.WaterQuality, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, ph, turbidity, dissolved_oxygen, temperature, conductivity, suspended_solids, timestamp FROM water_quality WHERE pool_id = ? ORDER BY timestamp DESC LIMIT ?", poolID, limit)
	if err != nil {
		return nil, fmt.Errorf("list water quality by pool: %w", err)
	}
	defer rows.Close()
	return scanWaterQualities(rows)
}

func scanWaterQualities(rows *sql.Rows) ([]model.WaterQuality, error) {
	var records []model.WaterQuality
	for rows.Next() {
		var wq model.WaterQuality
		var ts string
		if err := rows.Scan(&wq.ID, &wq.PoolID, &wq.PH, &wq.Turbidity, &wq.DissolvedOxygen, &wq.Temperature, &wq.Conductivity, &wq.SuspendedSolids, &ts); err != nil {
			return nil, fmt.Errorf("scan water quality row: %w", err)
		}
		wq.Timestamp, _ = time.Parse(time.RFC3339, ts)
		records = append(records, wq)
	}
	return records, nil
}
