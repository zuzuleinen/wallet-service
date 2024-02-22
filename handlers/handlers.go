package handlers

import (
	"encoding/json"
	"net/http"
)

// jsonResponse encodes resp to JSON and write it to w
func jsonResponse(w http.ResponseWriter, resp any) {
	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	if err := e.Encode(&resp); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// jsonSuccess returns a json success true with a msg and statusCode
func jsonSuccess(w http.ResponseWriter, msg string, statusCode int) {
	resp := struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}{Success: true, Msg: msg}

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := e.Encode(&resp); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// jsonError returns a json success false with a msg
func jsonError(w http.ResponseWriter, msg string, statusCode int) {
	resp := struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}{Success: false, Msg: msg}

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := e.Encode(&resp); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
