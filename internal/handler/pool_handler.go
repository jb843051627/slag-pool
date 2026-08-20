package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type PoolHandler struct {
	svc *service.PoolService
}

func NewPoolHandler(svc *service.PoolService) *PoolHandler {
	return &PoolHandler{svc: svc}
}

func (h *PoolHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/pools", h.Create)
	mux.HandleFunc("GET /api/pools", h.List)
	mux.HandleFunc("GET /api/pools/{id}", h.Get)
	mux.HandleFunc("PUT /api/pools/{id}", h.Update)
	mux.HandleFunc("DELETE /api/pools/{id}", h.Delete)
	mux.HandleFunc("GET /api/pools/{id}/status", h.GetStatus)
}

func (h *PoolHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pool model.SlagPool
	if err := json.NewDecoder(r.Body).Decode(&pool); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.svc.Create(r.Context(), &pool)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *PoolHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 100
	}
	pools, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pools)
}

func (h *PoolHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	pool, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, pool)
}

func (h *PoolHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var pool model.SlagPool
	if err := json.NewDecoder(r.Body).Decode(&pool); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	pool.ID = id
	if err := h.svc.Update(r.Context(), &pool); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pool)
}

func (h *PoolHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *PoolHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	status, err := h.svc.GetStatus(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, status)
}
