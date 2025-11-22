package http

import (
	"encoding/json"
	"net/http"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
)

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func respondError(w http.ResponseWriter, statusCode int, code, message string) {
	respondJSON(w, statusCode, dto.NewErrorResponse(code, message))
}
