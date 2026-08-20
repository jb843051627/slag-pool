package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type MaintenanceHandler struct {
	svc *service.MaintenanceService
}

func NewMaintenanceHandler(svc *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: svc}
}

func (h *MaintenanceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/maintenance", h.Create)
	mux.HandleFunc("GET /api/maintenance", h.List)
	mux.HandleFunc("GET /api/pools/{id}/maintenance", h.ListByPool)
	mux.HandleFunc("POST /api/maintenance/{id}/complete", h.Complete)
	mux.HandleFunc("POST /api/maintenance/schedule", h.Schedule)
}

func (h *MaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task model.MaintenanceTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Create(r.Context(), &task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *MaintenanceHandler) List(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.ParseInt(r.URL.Query().Get("pool_id"), 10, 64)
	if err != nil || poolID == 0 {
		writeError(w, http.StatusBadRequest, "pool_id query parameter required")
		return
	}
	tasks, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *MaintenanceHandler) ListByPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	tasks, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *MaintenanceHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Cost float64 `json:"cost"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Complete(r.Context(), id, req.Cost); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (h *MaintenanceHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	var task model.MaintenanceTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Schedule(r.Context(), &task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
