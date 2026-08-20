package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jb843051627/slag-pool/internal/store"
)

type ReportService struct {
	store *store.Store
}

func NewReportService(s *store.Store) *ReportService {
	return &ReportService{store: s}
}

func (s *ReportService) GeneratePoolReport(ctx context.Context, poolID int64) (map[string]interface{}, error) {
	pool, err := s.store.Pools().GetByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, store.ErrPoolNotFound
	}
	batches, _ := s.store.Batches().ListByPool(ctx, poolID)
	readings, _ := s.store.Readings().GetLatest(ctx, poolID)
	alerts, _ := s.store.Alerts().ListByPool(ctx, poolID)
	tasks, _ := s.store.Maintenance().ListByPool(ctx, poolID)
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].CreatedAt.After(alerts[j].CreatedAt)
	})

	activeAlerts := 0
	for _, a := range alerts {
		if a.Status == "active" {
			activeAlerts++
		}
	}
	activeBatches := 0
	for _, b := range batches {
		if b.Status == "cooling" {
			activeBatches++
		}
	}
	return map[string]interface{}{
		"pool":           pool,
		"latest_readings": readings,
		"batch_summary": map[string]interface{}{
			"total_batches":   len(batches),
			"active_batches":  activeBatches,
		},
		"alert_count":      len(alerts),
		"active_alerts":    activeAlerts,
		"maintenance_count": len(tasks),
	}, nil
}

func (s *ReportService) GenerateBatchReport(ctx context.Context, batchID int64) (map[string]interface{}, error) {
	batch, err := s.store.Batches().GetByID(ctx, batchID)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, store.ErrBatchNotFound
	}
	pool, _ := s.store.Pools().GetByID(ctx, batch.PoolID)
	readings, _ := s.store.Readings().ListByPool(ctx, batch.PoolID, 100)
	return map[string]interface{}{
		"batch":    batch,
		"pool":     pool,
		"readings": readings,
	}, nil
}

func (s *ReportService) ExportReadingsCSV(ctx context.Context, poolID int64, start, end time.Time) (string, error) {
	readings, err := s.store.Readings().ListByPool(ctx, poolID, 10000)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("id,sensor_id,pool_id,value,unit,quality,timestamp\n")
	for _, r := range readings {
		if !r.Timestamp.Before(start) && !r.Timestamp.After(end) {
			sb.WriteString(fmt.Sprintf("%d,%d,%d,%s,%s,%s,%s\n",
				r.ID, r.SensorID, r.PoolID,
				strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%f", r.Value), "0"), "."),
				r.Unit, r.Quality, r.Timestamp.Format(time.RFC3339)))
		}
	}
	return sb.String(), nil
}

func (s *ReportService) GenerateMaintenanceReport(ctx context.Context, poolID int64) (map[string]interface{}, error) {
	tasks, err := s.store.Maintenance().ListByPool(ctx, poolID)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ScheduledDate.After(tasks[j].ScheduledDate)
	})
	pool, _ := s.store.Pools().GetByID(ctx, poolID)
	pending, completed := 0, 0
	for _, t := range tasks {
		switch t.Status {
		case "pending", "scheduled", "in_progress":
			pending++
		case "completed":
			completed++
		}
	}
	return map[string]interface{}{
		"pool":      pool,
		"tasks":     tasks,
		"total":     len(tasks),
		"pending":   pending,
		"completed": completed,
	}, nil
}
