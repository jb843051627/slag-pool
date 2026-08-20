package model

import "time"

type Equipment struct {
	ID               int64     `json:"id"`
	PoolID           int64     `json:"pool_id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Model            string    `json:"model"`
	Manufacturer     string    `json:"manufacturer"`
	InstallDate      time.Time `json:"install_date"`
	Status           string    `json:"status"`
	LastMaintenance  *time.Time `json:"last_maintenance,omitempty"`
	NextMaintenance  *time.Time `json:"next_maintenance,omitempty"`
}

const (
	EquipTypePump       = "pump"
	EquipTypeValve      = "valve"
	EquipTypeHeatExch   = "heat_exchanger"
	EquipTypeFilter     = "filter"
	EquipTypeSensor     = "sensor"
	EquipTypeMotor      = "motor"

	EquipStatusRunning   = "running"
	EquipStatusIdle      = "idle"
	EquipStatusFaulty    = "faulty"
	EquipStatusMaint     = "maintenance"
	EquipStatusRetired   = "retired"
)
