package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type PoolStore struct {
	db    *sql.DB
	mu    sync.RWMutex
	cache map[int64]*model.SlagPool
}

func (s *PoolStore) GetByID(ctx context.Context, id int64) (*model.SlagPool, error) {
	s.mu.RLock()
	if p, ok := s.cache[id]; ok {
		s.mu.RUnlock()
		return p, nil
	}
	s.mu.RUnlock()
	row := s.db.QueryRowContext(ctx,
		"SELECT id, name, location, capacity, unit, status, max_temp, min_flow, created_at, updated_at FROM pools WHERE id = ?", id)
	var p model.SlagPool
	var createdAt, updatedAt string
	if err := row.Scan(&p.ID, &p.Name, &p.Location, &p.Capacity, &p.Unit, &p.Status, &p.MaxTemp, &p.MinFlow, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan pool: %w", err)
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	s.mu.Lock()
	s.cache[id] = &p
	s.mu.Unlock()
	return &p, nil
}

func (s *PoolStore) List(ctx context.Context, limit, offset int) ([]model.SlagPool, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, location, capacity, unit, status, max_temp, min_flow, created_at, updated_at FROM pools LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list pools: %w", err)
	}
	defer rows.Close()
	return scanPools(rows)
}

func (s *PoolStore) ListByStatus(ctx context.Context, status string) ([]model.SlagPool, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, location, capacity, unit, status, max_temp, min_flow, created_at, updated_at FROM pools WHERE status = ?", status)
	if err != nil {
		return nil, fmt.Errorf("list pools by status: %w", err)
	}
	defer rows.Close()
	return scanPools(rows)
}

func (s *PoolStore) Create(ctx context.Context, pool *model.SlagPool) (int64, error) {
	if pool == nil || pool.Name == "" {
		return 0, ErrInvalidInput
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO pools (name, location, capacity, unit, status, max_temp, min_flow, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?)",
		pool.Name, pool.Location, pool.Capacity, pool.Unit, pool.Status, pool.MaxTemp, pool.MinFlow,
		pool.CreatedAt.Format(time.RFC3339), pool.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create pool: %w", err)
	}
	id, _ := res.LastInsertId()
	pool.ID = id
	s.mu.Lock()
	s.cache[id] = pool
	s.mu.Unlock()
	return id, nil
}

func (s *PoolStore) Update(ctx context.Context, pool *model.SlagPool) error {
	if pool == nil || pool.ID == 0 {
		return ErrInvalidInput
	}
	_, err := s.db.ExecContext(ctx,
		"UPDATE pools SET name=?, location=?, capacity=?, unit=?, status=?, max_temp=?, min_flow=?, updated_at=? WHERE id=?",
		pool.Name, pool.Location, pool.Capacity, pool.Unit, pool.Status, pool.MaxTemp, pool.MinFlow,
		pool.UpdatedAt.Format(time.RFC3339), pool.ID)
	if err != nil {
		return fmt.Errorf("update pool: %w", err)
	}
	s.mu.Lock()
	s.cache[pool.ID] = pool
	s.mu.Unlock()
	return nil
}

func (s *PoolStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM pools WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete pool: %w", err)
	}
	s.mu.Lock()
	delete(s.cache, id)
	s.mu.Unlock()
	return nil
}

func scanPools(rows *sql.Rows) ([]model.SlagPool, error) {
	var pools []model.SlagPool
	for rows.Next() {
		var p model.SlagPool
		var createdAt, updatedAt string
		if err := rows.Scan(&p.ID, &p.Name, &p.Location, &p.Capacity, &p.Unit, &p.Status, &p.MaxTemp, &p.MinFlow, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan pool row: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		pools = append(pools, p)
	}
	return pools, nil
}
