package api

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func jsonRes(rw http.ResponseWriter, status int, message string, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	out, _ := json.Marshal(jsonResponse{Status: status, Message: message, Data: data})
	rw.Write(out)
}
