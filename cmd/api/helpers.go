package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	var js []byte
	var err error

	if app.config.env == "dev" {
		js, err = json.MarshalIndent(data, "", "\t")
	} else {
		js, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
