package store

import "errors"

var (
	ErrPoolNotFound       = errors.New("pool not found")
	ErrBatchNotFound      = errors.New("batch not found")
	ErrSensorNotFound     = errors.New("sensor not found")
	ErrAlertNotFound      = errors.New("alert not found")
	ErrMaintenanceNotFound = errors.New("maintenance task not found")
	ErrEquipmentNotFound  = errors.New("equipment not found")
	ErrQualityNotFound    = errors.New("water quality record not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("conflict")
)
