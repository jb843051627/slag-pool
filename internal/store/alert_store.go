package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type AlertStore struct {
	db    *sql.DB
	mu    sync.RWMutex
	cache map[int64][]model.Alert
}

func (s *AlertStore) Create(ctx context.Context, alert *model.Alert) (int64, error) {
	if alert == nil || alert.PoolID == 0 {
		return 0, ErrInvalidInput
	}
	var sensorID any
	if alert.SensorID != nil {
		sensorID = *alert.SensorID
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO alerts (pool_id, sensor_id, type, level, message, status, routed_to, created_at) VALUES (?,?,?,?,?,?,?,?)",
		alert.PoolID, sensorID, alert.Type, alert.Level, alert.Message, alert.Status, alert.RoutedTo, alert.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create alert: %w", err)
	}
	id, _ := res.LastInsertId()
	alert.ID = id
	s.mu.Lock()
	delete(s.cache, alert.PoolID)
	s.mu.Unlock()
	return id, nil
}

func (s *AlertStore) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, pool_id, sensor_id, type, level, message, status, routed_to, created_at, acknowledged_at, resolved_at FROM alerts WHERE id = ?", id)
	var a model.Alert
	var createdAt, ackAt, resAt string
	var sensorID sql.NullInt64
	if err := row.Scan(&a.ID, &a.PoolID, &sensorID, &a.Type, &a.Level, &a.Message, &a.Status, &a.RoutedTo, &createdAt, &ackAt, &resAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan alert: %w", err)
	}
	if sensorID.Valid {
		sid := sensorID.Int64
		a.SensorID = &sid
	}
	a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if ackAt != "" {
		t, _ := time.Parse(time.RFC3339, ackAt)
		a.AcknowledgedAt = &t
	}
	if resAt != "" {
		t, _ := time.Parse(time.RFC3339, resAt)
		a.ResolvedAt = &t
	}
	return &a, nil
}

func (s *AlertStore) ListByPool(ctx context.Context, poolID int64) ([]model.Alert, error) {
	s.mu.RLock()
	if cached, ok := s.cache[poolID]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, sensor_id, type, level, message, status, routed_to, created_at, acknowledged_at, resolved_at FROM alerts WHERE pool_id = ?", poolID)
	if err != nil {
		return nil, fmt.Errorf("list alerts by pool: %w", err)
	}
	defer rows.Close()
	alerts, err := scanAlerts(rows)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[poolID] = alerts
	s.mu.Unlock()
	return alerts, nil
}

func (s *AlertStore) ListByStatus(ctx context.Context, status string) ([]model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, sensor_id, type, level, message, status, routed_to, created_at, acknowledged_at, resolved_at FROM alerts WHERE status = ?", status)
	if err != nil {
		return nil, fmt.Errorf("list alerts by status: %w", err)
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (s *AlertStore) Update(ctx context.Context, alert *model.Alert) error {
	if alert == nil || alert.ID == 0 {
		return ErrInvalidInput
	}
	var sensorID any
	if alert.SensorID != nil {
		sensorID = *alert.SensorID
	}
	var ackVal, resVal any
	if alert.AcknowledgedAt != nil {
		ackVal = alert.AcknowledgedAt.Format(time.RFC3339)
	}
	if alert.ResolvedAt != nil {
		resVal = alert.ResolvedAt.Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx,
		"UPDATE alerts SET pool_id=?, sensor_id=?, type=?, level=?, message=?, status=?, routed_to=?, acknowledged_at=?, resolved_at=? WHERE id=?",
		alert.PoolID, sensorID, alert.Type, alert.Level, alert.Message, alert.Status, alert.RoutedTo, ackVal, resVal, alert.ID)
	if err != nil {
		return fmt.Errorf("update alert: %w", err)
	}
	s.mu.Lock()
	delete(s.cache, alert.PoolID)
	s.mu.Unlock()
	return nil
}

func (s *AlertStore) Acknowledge(ctx context.Context, id int64, routedTo string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE alerts SET status=?, routed_to=?, acknowledged_at=? WHERE id=?",
		model.AlertStatusAcknowledged, routedTo, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("acknowledge alert: %w", err)
	}
	s.mu.Lock()
	for k := range s.cache {
		delete(s.cache, k)
	}
	s.mu.Unlock()
	return nil
}

func scanAlerts(rows *sql.Rows) ([]model.Alert, error) {
	var alerts []model.Alert
	for rows.Next() {
		var a model.Alert
		var createdAt, ackAt, resAt string
		var sensorID sql.NullInt64
		if err := rows.Scan(&a.ID, &a.PoolID, &sensorID, &a.Type, &a.Level, &a.Message, &a.Status, &a.RoutedTo, &createdAt, &ackAt, &resAt); err != nil {
			return nil, fmt.Errorf("scan alert row: %w", err)
		}
		if sensorID.Valid {
			sid := sensorID.Int64
			a.SensorID = &sid
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		if ackAt != "" {
			t, _ := time.Parse(time.RFC3339, ackAt)
			a.AcknowledgedAt = &t
		}
		if resAt != "" {
			t, _ := time.Parse(time.RFC3339, resAt)
			a.ResolvedAt = &t
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}
