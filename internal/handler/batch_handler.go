package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type BatchHandler struct {
	svc *service.BatchService
}

func NewBatchHandler(svc *service.BatchService) *BatchHandler {
	return &BatchHandler{svc: svc}
}

func (h *BatchHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batches", h.Create)
	mux.HandleFunc("GET /api/batches", h.List)
	mux.HandleFunc("GET /api/pools/{id}/batches", h.ListByPool)
	mux.HandleFunc("POST /api/batches/{id}/start", h.Start)
	mux.HandleFunc("POST /api/batches/{id}/complete", h.Complete)
}

func (h *BatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var batch model.CoolingBatch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Create(r.Context(), &batch)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *BatchHandler) List(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.ParseInt(r.URL.Query().Get("pool_id"), 10, 64)
	if err != nil || poolID == 0 {
		writeError(w, http.StatusBadRequest, "pool_id query parameter required")
		return
	}
	batches, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, batches)
}

func (h *BatchHandler) ListByPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pool id")
		return
	}
	batches, err := h.svc.ListByPool(r.Context(), poolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, batches)
}

func (h *BatchHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Start(r.Context(), id); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (h *BatchHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		ActualTemp float64 `json:"actual_temp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Complete(r.Context(), id, req.ActualTemp); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}
