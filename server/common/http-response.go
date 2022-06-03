package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RespondWithErr(w http.ResponseWriter, statusCode int, err error) {
	m := make(map[string]string)

	if statusCode == http.StatusInternalServerError {
		m["message"] = "Internal server error"
	} else {
		m["message"] = err.Error()
	}

	log.Printf("Returning status %d reason: %s\n", statusCode, err)

	RespondWithJson(w, m, statusCode)
}

func RespondWithInternalErr(w http.ResponseWriter, err error) {
	RespondWithErr(w, http.StatusInternalServerError, err)
}

func RespondWithJson(w http.ResponseWriter, v interface{}, statusCode int) {
	content, e := json.Marshal(v)
	if e != nil {
		fmt.Printf("Could not send json response due to: %s\n", e)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(content)
}
