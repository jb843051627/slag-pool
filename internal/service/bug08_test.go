package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

func TestBug08_ReportSortPreservesStoreOrder(t *testing.T) {
	st := newTestStore(t)
	svc := NewReportService(st)
	ctx := context.Background()
	pid, err := st.Pools().Create(ctx, &model.SlagPool{Name: "p1", Status: model.PoolStatusActive})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if _, err := st.Alerts().Create(ctx, &model.Alert{PoolID: pid, Type: "temperature", Level: "high", Message: "m1", Status: model.AlertStatusActive, RoutedTo: "ops", CreatedAt: now.Add(-2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.Alerts().Create(ctx, &model.Alert{PoolID: pid, Type: "pressure", Level: "low", Message: "m2", Status: model.AlertStatusActive, RoutedTo: "ops", CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	initial, err := st.Alerts().ListByPool(ctx, pid)
	if err != nil {
		t.Fatal(err)
	}
	initialIDs := make([]int64, len(initial))
	for i := range initial {
		initialIDs[i] = initial[i].ID
	}
	if _, err := svc.GeneratePoolReport(ctx, pid); err != nil {
		t.Fatal(err)
	}
	after, err := st.Alerts().ListByPool(ctx, pid)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(initialIDs) {
		t.Fatalf("alert count changed: got %d want %d", len(after), len(initialIDs))
	}
	for i := range initialIDs {
		if after[i].ID != initialIDs[i] {
			t.Fatalf("store cache order polluted by report: index %d got id %d want id %d", i, after[i].ID, initialIDs[i])
		}
	}
}