package httphandler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bosskrub9992/fuel-management-backend/library/errs"
)

func sendJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	ctx := r.Context()
	body, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func sendBadRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := errs.ErrBadRequest
	body, err := json.Marshal(response)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	w.WriteHeader(response.Status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
