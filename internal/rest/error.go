package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteError(w http.ResponseWriter, log *slog.Logger, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	e := Error{
		Message: message,
	}

	encoder := json.NewEncoder(w)
	if er := encoder.Encode(e); er != nil {
		log.Error("Error writing response", "error", er)
	}
}
