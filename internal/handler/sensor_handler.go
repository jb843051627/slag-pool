package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type SensorHandler struct {
	svc *service.SensorService
}

func NewSensorHandler(svc *service.SensorService) *SensorHandler {
	return &SensorHandler{svc: svc}
}

func (h *SensorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sensors", h.Create)
	mux.HandleFunc("GET /api/sensors", h.List)
	mux.HandleFunc("GET /api/pools/{id}/sensors", h.ListByPool)
	mux.HandleFunc("PUT /api/sensors/{id}", h.Update)
	mux.HandleFunc("POST /api/sensors/{id}/assign", h.Assign)
}

func (h *SensorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sensor model.Sensor
	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Create(r.Context(), &sensor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *SensorHandler) List(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.ParseInt(r.URL.Query().Get("pool_id"), 10, 64)
	if err != nil || poolID == 0 {
		writeError(w, http.StatusBadRequest, "pool_id query parameter required")
		return
	}
	sensors, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sensors)
}

func (h *SensorHandler) ListByPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	sensors, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sensors)
}

func (h *SensorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var sensor model.Sensor
	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sensor.ID = id
	if err := h.svc.Update(r.Context(), &sensor); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sensor)
}

func (h *SensorHandler) Assign(w http.ResponseWriter, r *http.Request) {
	sensorID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid sensor id")
		return
	}
	var req struct {
		PoolID int64 `json:"pool_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.AssignToPool(r.Context(), sensorID, req.PoolID); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}
