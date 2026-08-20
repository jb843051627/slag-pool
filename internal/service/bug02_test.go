package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

func TestBug02_ListByPoolDoesNotReorderCache(t *testing.T) {
	st := newTestStore(t)
	svc := NewAlertService(st)
	ctx := context.Background()
	now := time.Now()
	if _, err := st.Alerts().Create(ctx, &model.Alert{PoolID: 1, Type: "temperature", Level: "high", Message: "m1", Status: model.AlertStatusActive, RoutedTo: "ops", CreatedAt: now.Add(-2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.Alerts().Create(ctx, &model.Alert{PoolID: 1, Type: "pressure", Level: "low", Message: "m2", Status: model.AlertStatusActive, RoutedTo: "ops", CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	initial, err := st.Alerts().ListByPool(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	initialIDs := make([]int64, len(initial))
	for i := range initial {
		initialIDs[i] = initial[i].ID
	}
	if _, err := svc.ListByPool(ctx, 1); err != nil {
		t.Fatal(err)
	}
	after, err := st.Alerts().ListByPool(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(initialIDs) {
		t.Fatalf("alert count changed: got %d want %d", len(after), len(initialIDs))
	}
	for i := range initialIDs {
		if after[i].ID != initialIDs[i] {
			t.Fatalf("store cache order polluted: index %d got id %d want id %d", i, after[i].ID, initialIDs[i])
		}
	}
}