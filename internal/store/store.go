package store

import (
	"database/sql"
	"fmt"

	"github.com/jb843051627/slag-pool/internal/model"
	_ "modernc.org/sqlite"
)

type Store struct {
	db          *sql.DB
	pool        *PoolStore
	batch       *BatchStore
	sensor      *SensorStore
	reading     *ReadingStore
	alert       *AlertStore
	maintenance *MaintenanceStore
	quality     *QualityStore
}

func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" || dbPath == ":memory:" {
		return nil, ErrInvalidInput
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	s := &Store{db: db}
	s.pool = &PoolStore{db: db, cache: make(map[int64]*model.SlagPool)}
	s.batch = &BatchStore{db: db, cache: make(map[int64][]model.CoolingBatch)}
	s.sensor = &SensorStore{db: db}
	s.reading = &ReadingStore{db: db}
	s.alert = &AlertStore{db: db, cache: make(map[int64][]model.Alert)}
	s.maintenance = &MaintenanceStore{db: db}
	s.quality = &QualityStore{db: db}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Pools() *PoolStore                 { return s.pool }
func (s *Store) Batches() *BatchStore              { return s.batch }
func (s *Store) Sensors() *SensorStore             { return s.sensor }
func (s *Store) Readings() *ReadingStore           { return s.reading }
func (s *Store) Alerts() *AlertStore               { return s.alert }
func (s *Store) Maintenance() *MaintenanceStore   { return s.maintenance }
func (s *Store) Quality() *QualityStore           { return s.quality }

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS pools (
		id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, location TEXT,
		capacity REAL, unit TEXT, status TEXT, max_temp REAL, min_flow REAL,
		created_at TEXT, updated_at TEXT
	);
	CREATE TABLE IF NOT EXISTS cooling_batches (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, batch_number TEXT,
		start_time TEXT, end_time TEXT, target_temp REAL, actual_temp REAL,
		volume REAL, source TEXT, status TEXT, created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS sensors (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, name TEXT,
		type TEXT, location TEXT, unit TEXT, min_value REAL, max_value REAL,
		status TEXT, created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS sensor_readings (
		id INTEGER PRIMARY KEY AUTOINCREMENT, sensor_id INTEGER, pool_id INTEGER,
		value REAL, unit TEXT, quality TEXT, timestamp TEXT
	);
	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, sensor_id INTEGER,
		type TEXT, level TEXT, message TEXT, status TEXT, routed_to TEXT,
		created_at TEXT, acknowledged_at TEXT, resolved_at TEXT
	);
	CREATE TABLE IF NOT EXISTS maintenance_tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, equipment_id INTEGER,
		type TEXT, description TEXT, status TEXT, priority TEXT, scheduled_date TEXT,
		completed_date TEXT, assigned_to TEXT, cost REAL, created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS equipment (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, name TEXT,
		type TEXT, model TEXT, manufacturer TEXT, install_date TEXT, status TEXT,
		last_maintenance TEXT, next_maintenance TEXT
	);
	CREATE TABLE IF NOT EXISTS water_quality (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, ph REAL,
		turbidity REAL, dissolved_oxygen REAL, temperature REAL, conductivity REAL,
		suspended_solids REAL, timestamp TEXT
	);
	CREATE TABLE IF NOT EXISTS alert_rules (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pool_id INTEGER, sensor_type TEXT,
		threshold REAL, operator TEXT, level TEXT, enabled INTEGER
	);
	`
	_, err := db.Exec(schema)
	return err
}
