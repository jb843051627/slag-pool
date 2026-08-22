package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

func TestBug07_BatchCompleteDetectsMissingTasks(t *testing.T) {
	st := newTestStore(t)
	svc := NewMaintenanceService(st, NewPoolService(st))
	ctx := context.Background()
	now := time.Now()
	id1, err := st.Maintenance().Create(ctx, &model.MaintenanceTask{PoolID: 1, Type: "inspection", Description: "d1", Status: model.MaintStatusPending, Priority: "normal", ScheduledDate: now, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	id2, err := st.Maintenance().Create(ctx, &model.MaintenanceTask{PoolID: 1, Type: "repair", Description: "d2", Status: model.MaintStatusPending, Priority: "normal", ScheduledDate: now, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.BatchComplete(ctx, []int64{id1, id2, 999}, 100)
	if err == nil {
		t.Fatal("BatchComplete: expected error for missing task")
	}
	t1, err := svc.Get(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if t1.Status == model.MaintStatusCompleted {
		t.Fatal("BatchComplete: task should not be completed when batch fails")
	}
}