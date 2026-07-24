package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
	}
}
