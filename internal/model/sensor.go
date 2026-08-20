package model

import "time"

type Sensor struct {
	ID        int64     `json:"id"`
	PoolID    int64     `json:"pool_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Location  string    `json:"location"`
	Unit      string    `json:"unit"`
	MinValue  float64   `json:"min_value"`
	MaxValue  float64   `json:"max_value"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	SensorTypeTemp      = "temperature"
	SensorTypeFlow       = "flow"
	SensorTypePressure   = "pressure"
	SensorTypePH         = "ph"
	SensorTypeTurbidity  = "turbidity"

	SensorStatusOnline   = "online"
	SensorStatusOffline  = "offline"
	SensorStatusFault    = "fault"
)
