package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type AlertHandler struct {
	svc *service.AlertService
}

func NewAlertHandler(svc *service.AlertService) *AlertHandler {
	return &AlertHandler{svc: svc}
}

func (h *AlertHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/alerts", h.List)
	mux.HandleFunc("GET /api/pools/{id}/alerts", h.ListByPool)
	mux.HandleFunc("POST /api/alerts", h.Create)
	mux.HandleFunc("POST /api/alerts/{id}/acknowledge", h.Acknowledge)
	mux.HandleFunc("POST /api/alerts/{id}/resolve", h.Resolve)
}

func (h *AlertHandler) List(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.ParseInt(r.URL.Query().Get("pool_id"), 10, 64)
	if err == nil && poolID > 0 {
		alerts, err := h.svc.ListByPool(r.Context(), poolID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, alerts)
		return
	}
	status := r.URL.Query().Get("status")
	if status != "" {
		alerts, err := h.svc.ListByStatus(r.Context(), status)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, alerts)
		return
	}
	writeError(w, http.StatusBadRequest, "pool_id or status query parameter required")
}

func (h *AlertHandler) ListByPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	alerts, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, alerts)
}

func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	var alert model.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Create(r.Context(), &alert)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *AlertHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		RoutedTo string `json:"routed_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Acknowledge(r.Context(), id, req.RoutedTo); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}

func (h *AlertHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Resolve(r.Context(), id); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
}
