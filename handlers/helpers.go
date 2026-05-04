package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/josesan28/proyecto-1-backend-web/models"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

func extractID(path string, position int) (int, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) <= position {
		return 0, strconv.ErrSyntax
	}
	return strconv.Atoi(parts[position])
}