package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, 200, envelope{"status": "healthy"}, nil)
}
