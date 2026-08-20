package service

import (
	"context"
	"time"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/store"
)

type QualityService struct {
	store *store.Store
}

func NewQualityService(s *store.Store) *QualityService {
	return &QualityService{store: s}
}

func (s *QualityService) Record(ctx context.Context, wq *model.WaterQuality) (int64, error) {
	if wq == nil || wq.PoolID == 0 {
		return 0, store.ErrInvalidInput
	}
	if wq.Timestamp.IsZero() {
		wq.Timestamp = time.Now()
	}
	return s.store.Quality().Create(ctx, wq)
}

func (s *QualityService) GetLatest(ctx context.Context, poolID int64) (*model.WaterQuality, error) {
	wq, err := s.store.Quality().GetLatest(ctx, poolID)
	if err != nil {
		return nil, err
	}
	if wq == nil {
		return nil, store.ErrQualityNotFound
	}
	return wq, nil
}

func (s *QualityService) Analyze(ctx context.Context, poolID int64) (*model.QualityReport, error) {
	wq, err := s.store.Quality().GetLatest(ctx, poolID)
	if err != nil {
		return nil, err
	}
	if wq == nil {
		return nil, store.ErrQualityNotFound
	}
	pool, _ := s.store.Pools().GetByID(ctx, poolID)
	report := &model.QualityReport{
		PoolID:      poolID,
		Latest:      *wq,
		Trend:       model.QualityTrendStable,
		Compliance:  true,
		GeneratedAt: time.Now(),
	}
	report.PoolName = pool.Name
	limits := model.QualityLimits
	if wq.PH < limits.PHMin || wq.PH > limits.PHMax {
		report.Compliance = false
		report.Issues = append(report.Issues, "pH out of range")
	}
	if wq.Turbidity > limits.TurbidityMax {
		report.Compliance = false
		report.Issues = append(report.Issues, "turbidity exceeds limit")
	}
	if wq.DissolvedOxygen < limits.DissolvedOxygenMin {
		report.Compliance = false
		report.Issues = append(report.Issues, "dissolved oxygen below minimum")
	}
	if wq.Conductivity > limits.ConductivityMax {
		report.Compliance = false
		report.Issues = append(report.Issues, "conductivity exceeds limit")
	}
	if wq.SuspendedSolids > limits.SuspendedSolidsMax {
		report.Compliance = false
		report.Issues = append(report.Issues, "suspended solids exceed limit")
	}
	return report, nil
}
