package utils

import (
	"net/http"
)

// InternalServerError : Send internal server error
func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}