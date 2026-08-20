package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jb843051627/slag-pool/internal/store"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (int64, error) {
	parts := strings.Split(r.URL.Path, "/")
	for _, p := range parts {
		if id, err := strconv.ParseInt(p, 10, 64); err == nil {
			return id, nil
		}
	}
	return 0, fmt.Errorf("invalid id in path")
}

func isNotFound(err error) bool {
	return errors.Is(err, store.ErrPoolNotFound) ||
		errors.Is(err, store.ErrBatchNotFound) ||
		errors.Is(err, store.ErrSensorNotFound) ||
		errors.Is(err, store.ErrAlertNotFound) ||
		errors.Is(err, store.ErrMaintenanceNotFound) ||
		errors.Is(err, store.ErrQualityNotFound)
}
