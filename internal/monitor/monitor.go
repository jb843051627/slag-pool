package monitor

import (
	"context"
	"log"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type Collector struct {
	pools    *service.PoolService
	sensors  *service.SensorService
	readings *service.ReadingService
	interval time.Duration
}

func NewCollector(ps *service.PoolService, ss *service.SensorService, rs *service.ReadingService) *Collector {
	return &Collector{pools: ps, sensors: ss, readings: rs, interval: 30 * time.Second}
}

func (c *Collector) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Collect(ctx); err != nil {
				log.Printf("collect error: %v", err)
			}
		}
	}
}

func (c *Collector) Collect(ctx context.Context) error {
	pools, err := c.pools.List(ctx, 100, 0)
	if err != nil {
		return err
	}
	for _, p := range pools {
		if p.Status != model.PoolStatusActive {
			continue
		}
		sensors, err := c.sensors.ListByPool(ctx, p.ID)
		if err != nil {
			continue
		}
		for _, sen := range sensors {
			if sen.Status != model.SensorStatusOnline {
				continue
			}
			reading := c.generateReading(&p, &sen)
			if _, err := c.readings.Ingest(ctx, reading); err != nil {
				log.Printf("ingest reading for sensor %d: %v", sen.ID, err)
			}
		}
	}
	return nil
}

func (c *Collector) generateReading(pool *model.SlagPool, sensor *model.Sensor) *model.SensorReading {
	now := time.Now()
	var value float64
	switch sensor.Type {
	case model.SensorTypeTemp:
		value = 45.0 + float64(now.Second())/10.0
	case model.SensorTypeFlow:
		value = pool.MinFlow + float64(now.Second())/5.0
	case model.SensorTypePressure:
		value = 1.5 + float64(now.Second())/20.0
	default:
		value = 7.0
	}
	return &model.SensorReading{
		SensorID:  sensor.ID,
		PoolID:    pool.ID,
		Value:     value,
		Unit:      sensor.Unit,
		Quality:   model.ReadingQualityGood,
		Timestamp: now,
	}
}

type Aggregator struct {
	readings *service.ReadingService
	interval time.Duration
}

func NewAggregator(rs *service.ReadingService) *Aggregator {
	return &Aggregator{readings: rs, interval: 60 * time.Second}
}

func (a *Aggregator) Run(ctx context.Context) {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.Aggregate(ctx); err != nil {
				log.Printf("aggregate error: %v", err)
			}
		}
	}
}

func (a *Aggregator) Aggregate(ctx context.Context) error {
	pools, err := a.readings.GetLatest(ctx, 0)
	_ = pools
	return err
}

type ThresholdChecker struct {
	pools    *service.PoolService
	alerts   *service.AlertService
	interval time.Duration
}

func NewThresholdChecker(ps *service.PoolService, as *service.AlertService) *ThresholdChecker {
	return &ThresholdChecker{pools: ps, alerts: as, interval: 30 * time.Second}
}

func (tc *ThresholdChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(tc.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := tc.Check(ctx); err != nil {
				log.Printf("threshold check error: %v", err)
			}
		}
	}
}

func (tc *ThresholdChecker) Check(ctx context.Context) error {
	pools, err := tc.pools.List(ctx, 100, 0)
	if err != nil {
		return err
	}
	for _, p := range pools {
		if p.Status != model.PoolStatusActive {
			continue
		}
		status, err := tc.pools.GetStatus(ctx, p.ID)
		if err != nil {
			continue
		}
		if p.MaxTemp > 0 && status.Temp > p.MaxTemp {
			alert := &model.Alert{
				PoolID:  p.ID,
				Type:    model.AlertTypeTempHigh,
				Level:   model.AlertLevelWarning,
				Message: "temperature exceeds threshold",
				Status:  model.AlertStatusActive,
			}
			tc.alerts.Create(ctx, alert)
		}
		if p.MinFlow > 0 && status.Flow < p.MinFlow {
			alert := &model.Alert{
				PoolID:  p.ID,
				Type:    model.AlertTypeFlowLow,
				Level:   model.AlertLevelWarning,
				Message: "flow below minimum threshold",
				Status:  model.AlertStatusActive,
			}
			tc.alerts.Create(ctx, alert)
		}
	}
	return nil
}
