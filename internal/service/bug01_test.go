package service

import (
	"context"
	"testing"

	"github.com/jb843051627/slag-pool/internal/store"
)

func TestBug01_GetMissingPoolReturnsError(t *testing.T) {
	st := newTestStore(t)
	svc := NewPoolService(st)
	ctx := context.Background()
	if _, err := svc.Get(ctx, 999); err != store.ErrPoolNotFound {
		t.Fatalf("Get missing pool: expected ErrPoolNotFound, got %v", err)
	}
	if _, err := svc.GetStatus(ctx, 999); err != store.ErrPoolNotFound {
		t.Fatalf("GetStatus missing pool: expected ErrPoolNotFound, got %v", err)
	}
}