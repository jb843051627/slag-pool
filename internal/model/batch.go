package model

import "time"

type CoolingBatch struct {
	ID          int64      `json:"id"`
	PoolID      int64      `json:"pool_id"`
	BatchNumber string     `json:"batch_number"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	TargetTemp  float64    `json:"target_temp"`
	ActualTemp  float64    `json:"actual_temp"`
	Volume      float64    `json:"volume"`
	Source      string     `json:"source"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

const (
	BatchStatusPending    = "pending"
	BatchStatusCooling    = "cooling"
	BatchStatusComplete   = "complete"
	BatchStatusAborted    = "aborted"
)

type BatchSummary struct {
	PoolID        int64   `json:"pool_id"`
	TotalBatches  int     `json:"total_batches"`
	ActiveBatches int     `json:"active_batches"`
	AvgCoolTime   float64 `json:"avg_cool_time_hours"`
	TotalVolume   float64 `json:"total_volume"`
}
