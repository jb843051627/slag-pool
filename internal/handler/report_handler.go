package handler

import (
	"net/http"
	"time"

	"github.com/jb843051627/slag-pool/internal/service"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/pool/{id}", h.PoolReport)
	mux.HandleFunc("GET /api/reports/batch/{id}", h.BatchReport)
	mux.HandleFunc("GET /api/reports/readings/{id}/csv", h.ReadingsCSV)
	mux.HandleFunc("GET /api/reports/maintenance/{id}", h.MaintenanceReport)
}

func (h *ReportHandler) PoolReport(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	report, err := h.svc.GeneratePoolReport(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) BatchReport(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch id")
		return
	}
	report, err := h.svc.GenerateBatchReport(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) ReadingsCSV(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	end := time.Now()
	start := end.AddDate(0, 0, -1)
	if s := r.URL.Query().Get("start"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			start = t
		}
	}
	if e := r.URL.Query().Get("end"); e != "" {
		if t, err := time.Parse(time.RFC3339, e); err == nil {
			end = t
		}
	}
	csv, err := h.svc.ExportReadingsCSV(r.Context(), id, start, end)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=readings.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}

func (h *ReportHandler) MaintenanceReport(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	report, err := h.svc.GenerateMaintenanceReport(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, report)
}
