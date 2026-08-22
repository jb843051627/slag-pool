package service

import (
	"context"
	"testing"
)

func TestBug09_CompleteMissingEntity(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	msvc := NewMaintenanceService(st, NewPoolService(st))
	if err := msvc.Complete(ctx, 999, 100); err == nil {
		t.Fatal("Maintenance Complete: expected error for missing task")
	}
	bsvc := NewBatchService(st, NewPoolService(st))
	if err := bsvc.Complete(ctx, 999, 25.0); err == nil {
		t.Fatal("Batch Complete: expected error for missing batch")
	}
}