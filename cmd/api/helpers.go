package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

// writeJson encoding data to json format and send it to client
func (app *application) writeJson(w http.ResponseWriter, r *http.Request, status int, data envelope, header http.Header) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// add any potential header
	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}
