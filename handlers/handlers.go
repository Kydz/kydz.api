package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func jsonizeResponse(data interface{}, w http.ResponseWriter) {
	jsonResponse, err := json.Marshal(data)
	var response string
	if err != nil {
		log.Fatal(err)
		errorResponse("Jsonize filed", w)
	} else {
		response = string(jsonResponse)
		normalResponse(response, w)
	}
}

func errorResponse(message string, w http.ResponseWriter) {
	var response = `{"error": true, "message":` + message + `}`
	normalResponse(response, w)
}

func normalResponse(response string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, string(response))
	if err != nil {
		log.Println("[Error]: normal response error: " + err.Error())
	}
}
