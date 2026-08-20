package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
)

func TestBug10_ExportCSVKeepsLocalTimestamps(t *testing.T) {
	st := newTestStore(t)
	svc := NewReportService(st)
	ctx := context.Background()
	pid, err := st.Pools().Create(ctx, &model.SlagPool{Name: "p1", Status: model.PoolStatusActive})
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Date(2026, 1, 1, 10, 30, 0, 0, time.Local)
	if _, err := st.Readings().Create(ctx, &model.SensorReading{SensorID: 1, PoolID: pid, Value: 25.5, Unit: "c", Quality: model.ReadingQualityGood, Timestamp: ts}); err != nil {
		t.Fatal(err)
	}
	csv, err := svc.ExportReadingsCSV(ctx, pid, ts.Add(-time.Hour), ts.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(csv, ts.Format(time.RFC3339)) {
		t.Fatalf("CSV missing local timestamp %s, got:\n%s", ts.Format(time.RFC3339), csv)
	}
}

func TestBug10_AlertKeepsLocalCreatedAt(t *testing.T) {
	st := newTestStore(t)
	svc := NewAlertService(st)
	ctx := context.Background()
	ts := time.Date(2026, 1, 1, 10, 30, 0, 0, time.Local)
	id, err := svc.Create(ctx, &model.Alert{PoolID: 1, Type: "temperature", Level: "high", Message: "m", Status: model.AlertStatusActive, CreatedAt: ts})
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.CreatedAt.Format(time.RFC3339) != ts.Format(time.RFC3339) {
		t.Fatalf("CreatedAt: got %v want %v", got.CreatedAt.Format(time.RFC3339), ts.Format(time.RFC3339))
	}
}