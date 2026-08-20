package model

import "time"

type WaterQuality struct {
	ID                int64     `json:"id"`
	PoolID            int64     `json:"pool_id"`
	PH                float64   `json:"ph"`
	Turbidity         float64   `json:"turbidity"`
	DissolvedOxygen   float64   `json:"dissolved_oxygen"`
	Temperature       float64   `json:"temperature"`
	Conductivity      float64   `json:"conductivity"`
	SuspendedSolids   float64   `json:"suspended_solids"`
	Timestamp         time.Time `json:"timestamp"`
}

type QualityReport struct {
	PoolID      int64     `json:"pool_id"`
	PoolName    string    `json:"pool_name"`
	Latest      WaterQuality `json:"latest"`
	Trend       string    `json:"trend"`
	Compliance  bool      `json:"compliance"`
	Issues      []string  `json:"issues"`
	GeneratedAt time.Time `json:"generated_at"`
}

const (
	QualityTrendImproving  = "improving"
	QualityTrendStable    = "stable"
	QualityTrendDegrading  = "degrading"
)

var QualityLimits = struct {
	PHMin, PHMax           float64
	TurbidityMax            float64
	DissolvedOxygenMin      float64
	ConductivityMax          float64
	SuspendedSolidsMax       float64
}{
	PHMin:              6.5,
	PHMax:              9.0,
	TurbidityMax:        50.0,
	DissolvedOxygenMin:  4.0,
	ConductivityMax:     2000.0,
	SuspendedSolidsMax:  100.0,
}
