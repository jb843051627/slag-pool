package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jb843051627/slag-pool/internal/alert"
	"github.com/jb843051627/slag-pool/internal/handler"
	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/monitor"
	"github.com/jb843051627/slag-pool/internal/service"
	"github.com/jb843051627/slag-pool/internal/store"
)

func main() {
	st, err := store.NewStore("slag_pool.db")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	poolSvc := service.NewPoolService(st)
	batchSvc := service.NewBatchService(st, poolSvc)
	sensorSvc := service.NewSensorService(st, poolSvc)
	readingSvc := service.NewReadingService(st, sensorSvc)
	alertSvc := service.NewAlertService(st)
	maintSvc := service.NewMaintenanceService(st, poolSvc)
	reportSvc := service.NewReportService(st)
	qualitySvc := service.NewQualityService(st)
	_ = qualitySvc

	seedData(poolSvc, sensorSvc)

	collector := monitor.NewCollector(poolSvc, sensorSvc, readingSvc)
	aggregator := monitor.NewAggregator(readingSvc)
	checker := monitor.NewThresholdChecker(poolSvc, alertSvc)
	engine := alert.NewEngine(alertSvc)
	notifier := alert.NewNotifier()
	router := alert.NewRouter(alertSvc)
	_ = notifier
	_ = router

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collector.Run(ctx)
	go aggregator.Run(ctx)
	go checker.Run(ctx)
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := engine.Evaluate(ctx); err != nil {
					log.Printf("alert engine: %v", err)
				}
			}
		}
	}()

	mux := http.NewServeMux()
	handler.NewPoolHandler(poolSvc).RegisterRoutes(mux)
	handler.NewBatchHandler(batchSvc).RegisterRoutes(mux)
	handler.NewSensorHandler(sensorSvc).RegisterRoutes(mux)
	handler.NewReadingHandler(readingSvc).RegisterRoutes(mux)
	handler.NewAlertHandler(alertSvc).RegisterRoutes(mux)
	handler.NewMaintenanceHandler(maintSvc).RegisterRoutes(mux)
	handler.NewReportHandler(reportSvc).RegisterRoutes(mux)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/templates/index.html")
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Println("server starting on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}

func seedData(poolSvc *service.PoolService, sensorSvc *service.SensorService) {
	pools, _ := poolSvc.List(context.Background(), 1, 0)
	if len(pools) > 0 {
		return
	}
	log.Println("seeding initial data...")
	seedPools := []model.SlagPool{
		{Name: "1号渣池", Location: "A区", Capacity: 500, Unit: "m³", Status: model.PoolStatusActive, MaxTemp: 80, MinFlow: 20},
		{Name: "2号渣池", Location: "B区", Capacity: 600, Unit: "m³", Status: model.PoolStatusActive, MaxTemp: 75, MinFlow: 25},
		{Name: "3号渣池", Location: "C区", Capacity: 450, Unit: "m³", Status: model.PoolStatusStandby, MaxTemp: 70, MinFlow: 15},
	}
	for _, p := range seedPools {
		pid, _ := poolSvc.Create(context.Background(), &p)
		seedSensors := []model.Sensor{
			{PoolID: pid, Name: p.Name + "-温度", Type: model.SensorTypeTemp, Unit: "°C", MinValue: 0, MaxValue: 150, Status: model.SensorStatusOnline},
			{PoolID: pid, Name: p.Name + "-流量", Type: model.SensorTypeFlow, Unit: "L/min", MinValue: 0, MaxValue: 500, Status: model.SensorStatusOnline},
			{PoolID: pid, Name: p.Name + "-压力", Type: model.SensorTypePressure, Unit: "kPa", MinValue: 0, MaxValue: 500, Status: model.SensorStatusOnline},
		}
		for _, s := range seedSensors {
			sensorSvc.Create(context.Background(), &s)
		}
	}
}
