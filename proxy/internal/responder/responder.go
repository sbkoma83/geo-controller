package responder

import (
	"encoding/json"
	"net/http"
)

// Responder реализует интерфейс Responder
type Responder struct{}

// NewResponder создает новый Responder
func NewResponder() *Responder {
	return &Responder{}
}

// OutputJSON отправляет данные в формате JSON
func (r *Responder) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if responseData != nil {
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
	}
}

// ErrorUnauthorized отправляет ответ с ошибкой 401 Unauthorized
func (r *Responder) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusUnauthorized, err)
}

// ErrorBadRequest отправляет ответ с ошибкой 400 Bad Request
func (r *Responder) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusBadRequest, err)
}

// ErrorForbidden отправляет ответ с ошибкой 403 Forbidden
func (r *Responder) ErrorForbidden(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusForbidden, err)
}

// ErrorInternal отправляет ответ с ошибкой 500 Internal Server Error
func (r *Responder) ErrorInternal(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusInternalServerError, err)
}

// sendError общий метод для отправки ошибок
func (r *Responder) sendError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	if err != nil {
		r.OutputJSON(w, map[string]string{"error": err.Error()})
	} else {
		r.OutputJSON(w, map[string]string{"error": "error"})
	}
}
