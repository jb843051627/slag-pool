package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type ReadingStore struct {
	db *sql.DB
}

func (s *ReadingStore) Create(ctx context.Context, reading *model.SensorReading) (int64, error) {
	if reading == nil || reading.SensorID == 0 {
		return 0, ErrInvalidInput
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO sensor_readings (sensor_id, pool_id, value, unit, quality, timestamp) VALUES (?,?,?,?,?,?)",
		reading.SensorID, reading.PoolID, reading.Value, reading.Unit, reading.Quality, reading.Timestamp.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create reading: %w", err)
	}
	id, _ := res.LastInsertId()
	reading.ID = id
	return id, nil
}

func (s *ReadingStore) ListBySensor(ctx context.Context, sensorID, limit int64) ([]model.SensorReading, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, sensor_id, pool_id, value, unit, quality, timestamp FROM sensor_readings WHERE sensor_id = ? ORDER BY timestamp DESC LIMIT ?", sensorID, limit)
	if err != nil {
		return nil, fmt.Errorf("list readings by sensor: %w", err)
	}
	defer rows.Close()
	return scanReadings(rows)
}

func (s *ReadingStore) ListByPool(ctx context.Context, poolID, limit int64) ([]model.SensorReading, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, sensor_id, pool_id, value, unit, quality, timestamp FROM sensor_readings WHERE pool_id = ? ORDER BY timestamp DESC LIMIT ?", poolID, limit)
	if err != nil {
		return nil, fmt.Errorf("list readings by pool: %w", err)
	}
	defer rows.Close()
	return scanReadings(rows)
}

func (s *ReadingStore) BatchCreate(ctx context.Context, readings []model.SensorReading) error {
	if len(readings) == 0 {
		return ErrInvalidInput
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO sensor_readings (sensor_id, pool_id, value, unit, quality, timestamp) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()
	for _, r := range readings {
		if _, err := stmt.ExecContext(ctx, r.SensorID, r.PoolID, r.Value, r.Unit, r.Quality, r.Timestamp.Format(time.RFC3339)); err != nil {
			return fmt.Errorf("exec reading insert: %w", err)
		}
	}
	return tx.Commit()
}

func (s *ReadingStore) GetLatest(ctx context.Context, poolID int64) ([]model.SensorReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT r.id, r.sensor_id, r.pool_id, r.value, r.unit, r.quality, r.timestamp
		FROM sensor_readings r
		INNER JOIN (
			SELECT sensor_id, MAX(timestamp) AS max_ts FROM sensor_readings WHERE pool_id = ? GROUP BY sensor_id
		) latest ON r.sensor_id = latest.sensor_id AND r.timestamp = latest.max_ts
		WHERE r.pool_id = ?`, poolID, poolID)
	if err != nil {
		return nil, fmt.Errorf("get latest readings: %w", err)
	}
	defer rows.Close()
	return scanReadings(rows)
}

func scanReadings(rows *sql.Rows) ([]model.SensorReading, error) {
	var readings []model.SensorReading
	for rows.Next() {
		var r model.SensorReading
		var ts string
		if err := rows.Scan(&r.ID, &r.SensorID, &r.PoolID, &r.Value, &r.Unit, &r.Quality, &ts); err != nil {
			return nil, fmt.Errorf("scan reading row: %w", err)
		}
		r.Timestamp, _ = time.Parse(time.RFC3339, ts)
		readings = append(readings, r)
	}
	return readings, nil
}
