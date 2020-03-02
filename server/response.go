package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func successResponse(w http.ResponseWriter, l *log.Logger, result interface{}) {
	data := map[string]interface{}{
		"response": result,
	}

	jsonData, jErr := json.Marshal(data)
	if jErr != nil {
		errorResponse(jErr, w, l, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func errorResponse(err error, w http.ResponseWriter, l *log.Logger, request int) {
	if request != http.StatusNotFound {
		l.Print(err)
	}

	data := map[string]interface{}{
		"response": false,
		"error":    err.Error(),
	}

	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(request)
	w.Write(jsonData)
}
