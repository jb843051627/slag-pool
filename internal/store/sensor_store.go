package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type SensorStore struct {
	db *sql.DB
}

func (s *SensorStore) GetByID(ctx context.Context, id int64) (*model.Sensor, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, pool_id, name, type, location, unit, min_value, max_value, status, created_at FROM sensors WHERE id = ?", id)
	var sen model.Sensor
	var createdAt string
	if err := row.Scan(&sen.ID, &sen.PoolID, &sen.Name, &sen.Type, &sen.Location, &sen.Unit, &sen.MinValue, &sen.MaxValue, &sen.Status, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan sensor: %w", err)
	}
	sen.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &sen, nil
}

func (s *SensorStore) ListByPool(ctx context.Context, poolID int64) ([]model.Sensor, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, name, type, location, unit, min_value, max_value, status, created_at FROM sensors WHERE pool_id = ?", poolID)
	if err != nil {
		return nil, fmt.Errorf("list sensors by pool: %w", err)
	}
	defer rows.Close()
	return scanSensors(rows)
}

func (s *SensorStore) ListByType(ctx context.Context, sensorType string) ([]model.Sensor, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, name, type, location, unit, min_value, max_value, status, created_at FROM sensors WHERE type = ?", sensorType)
	if err != nil {
		return nil, fmt.Errorf("list sensors by type: %w", err)
	}
	defer rows.Close()
	return scanSensors(rows)
}

func (s *SensorStore) Create(ctx context.Context, sensor *model.Sensor) (int64, error) {
	if sensor == nil || sensor.PoolID == 0 || sensor.Name == "" {
		return 0, ErrInvalidInput
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO sensors (pool_id, name, type, location, unit, min_value, max_value, status, created_at) VALUES (?,?,?,?,?,?,?,?,?)",
		sensor.PoolID, sensor.Name, sensor.Type, sensor.Location, sensor.Unit,
		sensor.MinValue, sensor.MaxValue, sensor.Status, sensor.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create sensor: %w", err)
	}
	id, _ := res.LastInsertId()
	sensor.ID = id
	return id, nil
}

func (s *SensorStore) Update(ctx context.Context, sensor *model.Sensor) error {
	if sensor == nil || sensor.ID == 0 {
		return ErrInvalidInput
	}
	_, err := s.db.ExecContext(ctx,
		"UPDATE sensors SET pool_id=?, name=?, type=?, location=?, unit=?, min_value=?, max_value=?, status=? WHERE id=?",
		sensor.PoolID, sensor.Name, sensor.Type, sensor.Location, sensor.Unit,
		sensor.MinValue, sensor.MaxValue, sensor.Status, sensor.ID)
	if err != nil {
		return fmt.Errorf("update sensor: %w", err)
	}
	return nil
}

func scanSensors(rows *sql.Rows) ([]model.Sensor, error) {
	var sensors []model.Sensor
	for rows.Next() {
		var sen model.Sensor
		var createdAt string
		if err := rows.Scan(&sen.ID, &sen.PoolID, &sen.Name, &sen.Type, &sen.Location, &sen.Unit, &sen.MinValue, &sen.MaxValue, &sen.Status, &createdAt); err != nil {
			return nil, fmt.Errorf("scan sensor row: %w", err)
		}
		sen.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		sensors = append(sensors, sen)
	}
	return sensors, nil
}
