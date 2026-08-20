package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

type MaintenanceStore struct {
	db *sql.DB
}

func (s *MaintenanceStore) Create(ctx context.Context, task *model.MaintenanceTask) (int64, error) {
	if task == nil || task.PoolID == 0 {
		return 0, ErrInvalidInput
	}
	var completedVal any
	if task.CompletedDate != nil {
		completedVal = task.CompletedDate.Format(time.RFC3339)
	}
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO maintenance_tasks (pool_id, equipment_id, type, description, status, priority, scheduled_date, completed_date, assigned_to, cost, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		task.PoolID, task.EquipmentID, task.Type, task.Description, task.Status, task.Priority,
		task.ScheduledDate.Format(time.RFC3339), completedVal, task.AssignedTo, task.Cost, task.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("create maintenance task: %w", err)
	}
	id, _ := res.LastInsertId()
	task.ID = id
	return id, nil
}

func (s *MaintenanceStore) GetByID(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, pool_id, equipment_id, type, description, status, priority, scheduled_date, completed_date, assigned_to, cost, created_at FROM maintenance_tasks WHERE id = ?", id)
	var t model.MaintenanceTask
	var scheduledDate, createdAt string
	var completedDate sql.NullString
	if err := row.Scan(&t.ID, &t.PoolID, &t.EquipmentID, &t.Type, &t.Description, &t.Status, &t.Priority, &scheduledDate, &completedDate, &t.AssignedTo, &t.Cost, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan maintenance task: %w", err)
	}
	t.ScheduledDate, _ = time.Parse(time.RFC3339, scheduledDate)
	if completedDate.Valid {
		cd, _ := time.Parse(time.RFC3339, completedDate.String)
		t.CompletedDate = &cd
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &t, nil
}

func (s *MaintenanceStore) ListByPool(ctx context.Context, poolID int64) ([]model.MaintenanceTask, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, pool_id, equipment_id, type, description, status, priority, scheduled_date, completed_date, assigned_to, cost, created_at FROM maintenance_tasks WHERE pool_id = ? ORDER BY scheduled_date DESC", poolID)
	if err != nil {
		return nil, fmt.Errorf("list maintenance tasks by pool: %w", err)
	}
	defer rows.Close()
	return scanMaintenanceTasks(rows)
}

func (s *MaintenanceStore) Update(ctx context.Context, task *model.MaintenanceTask) error {
	if task == nil || task.ID == 0 {
		return ErrInvalidInput
	}
	var completedVal any
	if task.CompletedDate != nil {
		completedVal = task.CompletedDate.Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx,
		"UPDATE maintenance_tasks SET pool_id=?, equipment_id=?, type=?, description=?, status=?, priority=?, scheduled_date=?, completed_date=?, assigned_to=?, cost=? WHERE id=?",
		task.PoolID, task.EquipmentID, task.Type, task.Description, task.Status, task.Priority,
		task.ScheduledDate.Format(time.RFC3339), completedVal, task.AssignedTo, task.Cost, task.ID)
	if err != nil {
		return fmt.Errorf("update maintenance task: %w", err)
	}
	return nil
}

func scanMaintenanceTasks(rows *sql.Rows) ([]model.MaintenanceTask, error) {
	var tasks []model.MaintenanceTask
	for rows.Next() {
		var t model.MaintenanceTask
		var scheduledDate, createdAt string
		var completedDate sql.NullString
		if err := rows.Scan(&t.ID, &t.PoolID, &t.EquipmentID, &t.Type, &t.Description, &t.Status, &t.Priority, &scheduledDate, &completedDate, &t.AssignedTo, &t.Cost, &createdAt); err != nil {
			return nil, fmt.Errorf("scan maintenance task row: %w", err)
		}
		t.ScheduledDate, _ = time.Parse(time.RFC3339, scheduledDate)
		if completedDate.Valid {
			cd, _ := time.Parse(time.RFC3339, completedDate.String)
			t.CompletedDate = &cd
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		tasks = append(tasks, t)
	}
	return tasks, nil
}
