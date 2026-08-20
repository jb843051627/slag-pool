package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type SensorService struct {
	store *store.Store
	pools *PoolService
}

func NewSensorService(s *store.Store, ps *PoolService) *SensorService {
	return &SensorService{store: s, pools: ps}
}

func (s *SensorService) Get(ctx context.Context, id int64) (*model.Sensor, error) {
	sen, err := s.store.Sensors().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sen == nil {
		return nil, store.ErrSensorNotFound
	}
	return sen, nil
}

func (s *SensorService) ListByPool(ctx context.Context, poolID int64) ([]model.Sensor, error) {
	return s.store.Sensors().ListByPool(ctx, poolID)
}

func (s *SensorService) Create(ctx context.Context, sensor *model.Sensor) (int64, error) {
	if sensor == nil || sensor.PoolID == 0 || sensor.Name == "" {
		return 0, store.ErrInvalidInput
	}
	if sensor.CreatedAt.IsZero() {
		sensor.CreatedAt = time.Now()
	}
	if sensor.Status == "" {
		sensor.Status = model.SensorStatusOnline
	}
	return s.store.Sensors().Create(ctx, sensor)
}

func (s *SensorService) Update(ctx context.Context, sensor *model.Sensor) error {
	if sensor == nil || sensor.ID == 0 {
		return store.ErrInvalidInput
	}
	return s.store.Sensors().Update(ctx, sensor)
}

func (s *SensorService) AssignToPool(ctx context.Context, sensorID, poolID int64) error {
	sen, err := s.Get(ctx, sensorID)
	if err != nil {
		return err
	}
	sen.PoolID = poolID
	return s.store.Sensors().Update(ctx, sen)
}
