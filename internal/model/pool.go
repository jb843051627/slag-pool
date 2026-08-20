package model

import "time"

type SlagPool struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Capacity  float64   `json:"capacity"`
	Unit      string    `json:"unit"`
	Status    string    `json:"status"`
	MaxTemp   float64   `json:"max_temp"`
	MinFlow   float64   `json:"min_flow"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	PoolStatusActive    = "active"
	PoolStatusStandby  = "standby"
	PoolStatusDraining  = "draining"
	PoolStatusMaint     = "maintenance"
	PoolStatusFault     = "fault"
)

type PoolStatus struct {
	PoolID     int64   `json:"pool_id"`
	Temp       float64 `json:"temp"`
	Flow       float64 `json:"flow"`
	Pressure   float64 `json:"pressure"`
	ActiveBatch string `json:"active_batch"`
	AlertCount int     `json:"alert_count"`
}
