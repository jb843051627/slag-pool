package service

import (
	"context"
	"testing"
)

func TestBug06_MissingEntityUpdate(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	ssvc := NewSensorService(st, NewPoolService(st))
	if err := ssvc.AssignToPool(ctx, 999, 1); err == nil {
		t.Fatal("AssignToPool: expected error for missing sensor")
	}
	asvc := NewAlertService(st)
	if err := asvc.Resolve(ctx, 999); err == nil {
		t.Fatal("Resolve: expected error for missing alert")
	}
}