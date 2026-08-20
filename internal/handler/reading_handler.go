package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type ReadingHandler struct {
	svc *service.ReadingService
}

func NewReadingHandler(svc *service.ReadingService) *ReadingHandler {
	return &ReadingHandler{svc: svc}
}

func (h *ReadingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/readings", h.Ingest)
	mux.HandleFunc("GET /api/sensors/{id}/readings", h.ListBySensor)
	mux.HandleFunc("GET /api/pools/{id}/readings", h.ListByPool)
	mux.HandleFunc("POST /api/readings/batch", h.BatchIngest)
	mux.HandleFunc("GET /api/pools/{id}/latest", h.GetLatest)
}

func (h *ReadingHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var reading model.SensorReading
	if err := json.NewDecoder(r.Body).Decode(&reading); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Ingest(r.Context(), &reading)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *ReadingHandler) ListBySensor(w http.ResponseWriter, r *http.Request) {
	sensorID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid sensor id")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}
	readings, err := h.svc.ListBySensor(r.Context(), sensorID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, readings)
}

func (h *ReadingHandler) ListByPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}
	readings, err := h.svc.ListByPool(r.Context(), poolID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, readings)
}

func (h *ReadingHandler) BatchIngest(w http.ResponseWriter, r *http.Request) {
	var batch model.ReadingBatch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.BatchIngest(r.Context(), &batch); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int{"count": len(batch.Readings)})
}

func (h *ReadingHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	readings, err := h.svc.GetLatest(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, readings)
}
