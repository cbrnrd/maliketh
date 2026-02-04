package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Msg     string `json:"msg,omitempty"` // Error message (Python client compatibility)
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, Response{Status: true, Message: message})
}

func SuccessData(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{Status: false, Msg: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
