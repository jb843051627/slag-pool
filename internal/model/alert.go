package model

import "time"

type Alert struct {
	ID         int64      `json:"id"`
	PoolID     int64      `json:"pool_id"`
	SensorID   *int64     `json:"sensor_id,omitempty"`
	Type       string     `json:"type"`
	Level      string     `json:"level"`
	Message    string     `json:"message"`
	Status     string     `json:"status"`
	RoutedTo   string     `json:"routed_to"`
	CreatedAt  time.Time  `json:"created_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

const (
	AlertTypeTempHigh    = "temperature_high"
	AlertTypeTempLow     = "temperature_low"
	AlertTypeFlowLow     = "flow_low"
	AlertTypePressure    = "pressure_abnormal"
	AlertTypePHAbnormal  = "ph_abnormal"
	AlertTypeSensorOffline = "sensor_offline"
	AlertTypeEquipFail   = "equipment_failure"

	AlertLevelInfo     = "info"
	AlertLevelWarning  = "warning"
	AlertLevelCritical = "critical"
	AlertLevelEmergency = "emergency"

	AlertStatusActive       = "active"
	AlertStatusAcknowledged = "acknowledged"
	AlertStatusResolved     = "resolved"
)

type AlertRule struct {
	ID        int64   `json:"id"`
	PoolID    int64   `json:"pool_id"`
	SensorType string `json:"sensor_type"`
	Threshold float64 `json:"threshold"`
	Operator  string `json:"operator"`
	Level     string `json:"level"`
	Enabled   bool    `json:"enabled"`
}
