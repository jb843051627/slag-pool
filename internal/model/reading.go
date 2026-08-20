package model

import "time"

type SensorReading struct {
	ID         int64     `json:"id"`
	SensorID   int64     `json:"sensor_id"`
	PoolID     int64     `json:"pool_id"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Quality    string    `json:"quality"`
	Timestamp  time.Time `json:"timestamp"`
}

const (
	ReadingQualityGood   = "good"
	ReadingQualityFair   = "fair"
	ReadingQualityPoor   = "poor"
	ReadingQualityStale  = "stale"
)

type ReadingBatch struct {
	PoolID   int64           `json:"pool_id"`
	Readings []SensorReading  `json:"readings"`
}

type AggregatedReading struct {
	PoolID    int64     `json:"pool_id"`
	SensorType string   `json:"sensor_type"`
	AvgValue  float64   `json:"avg_value"`
	MinValue  float64   `json:"min_value"`
	MaxValue  float64   `json:"max_value"`
	Count     int       `json:"count"`
	Timestamp time.Time `json:"timestamp"`
}
