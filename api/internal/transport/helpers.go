package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Error  string `json:"error"`
}

type loggingWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func jsonResp(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errResp(w http.ResponseWriter, status int, code, msg string) {
	jsonResp(w, status, ErrorResponse{Status: status, Code: code, Error: msg})
}

func badRequest(w http.ResponseWriter, msg string) {
	errResp(w, http.StatusBadRequest, "BAD_REQUEST", msg)
}

func notFound(w http.ResponseWriter, msg string) {
	errResp(w, http.StatusNotFound, "NOT_FOUND", msg)
}

func internalError(w http.ResponseWriter, msg string) {
	errResp(w, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

func unauthorized(w http.ResponseWriter) {
	errResp(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
}

func conflict(w http.ResponseWriter, msg string) {
	errResp(w, http.StatusConflict, "CONFLICT", msg)
}

func mapStoreError(err error) string {
	if err == nil {
		return ""
	}
	log.Printf("[STORE_ERROR] %v (type: %T)", err, err)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			if strings.Contains(pgErr.Message, "users_username_key") {
				return "username already exists"
			}
			if strings.Contains(pgErr.Message, "users_email_key") {
				return "email already exists"
			}
			return "duplicate entry"
		case "23503":
			return "referenced record not found"
		case "23514":
			return "invalid data"
		}
	}
	return "internal server error"
}