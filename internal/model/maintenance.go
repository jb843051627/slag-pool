package model

import "time"

type MaintenanceTask struct {
	ID            int64      `json:"id"`
	PoolID        int64      `json:"pool_id"`
	EquipmentID   int64      `json:"equipment_id"`
	Type          string     `json:"type"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Priority      string     `json:"priority"`
	ScheduledDate time.Time  `json:"scheduled_date"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`
	AssignedTo    string     `json:"assigned_to"`
	Cost          float64    `json:"cost"`
	CreatedAt     time.Time  `json:"created_at"`
}

const (
	MaintTypeInspection  = "inspection"
	MaintTypeRepair      = "repair"
	MaintTypeReplace     = "replace"
	MaintTypeCalibrate   = "calibrate"
	MaintTypeClean      = "clean"

	MaintStatusPending   = "pending"
	MaintStatusScheduled = "scheduled"
	MaintStatusInProgress = "in_progress"
	MaintStatusCompleted = "completed"
	MaintStatusCancelled = "cancelled"

	MaintPriorityLow    = "low"
	MaintPriorityMedium = "medium"
	MaintPriorityHigh   = "high"
	MaintPriorityUrgent = "urgent"
)

type MaintenanceSchedule struct {
	PoolID    int64              `json:"pool_id"`
	Tasks     []MaintenanceTask  `json:"tasks"`
	NextDue   *time.Time          `json:"next_due,omitempty"`
}
