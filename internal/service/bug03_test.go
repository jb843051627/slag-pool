package service

import (
	"context"
	"sync"
	"testing"
)

func TestBug03_GetLatestCacheConcurrentSafe(t *testing.T) {
	st := newTestStore(t)
	svc := NewReadingService(st, NewSensorService(st, NewPoolService(st)))
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := svc.GetLatest(ctx, 1); err != nil {
				t.Errorf("GetLatest: %v", err)
			}
		}()
	}
	wg.Wait()
}