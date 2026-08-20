package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/slag-pool/internal/store"
)

func TestBug04_NotFoundErrorsPreserveSentinel(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	asvc := NewAlertService(st)
	if _, err := asvc.Get(ctx, 999); !errors.Is(err, store.ErrAlertNotFound) {
		t.Fatalf("Alert Get: expected ErrAlertNotFound via errors.Is, got %v", err)
	}
	msvc := NewMaintenanceService(st, NewPoolService(st))
	if _, err := msvc.Get(ctx, 999); !errors.Is(err, store.ErrMaintenanceNotFound) {
		t.Fatalf("Maintenance Get: expected ErrMaintenanceNotFound via errors.Is, got %v", err)
	}
}