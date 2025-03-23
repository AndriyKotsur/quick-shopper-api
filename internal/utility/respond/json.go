package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

type Meta struct {
	Size  int `json:"size"`
	Total int `json:"total"`
}

type Standard struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}

	if string(data) == "null" {
		_, _ = w.Write([]byte("[]"))
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
}
