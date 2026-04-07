package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/cars", app.createCarHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/cars/:id", app.updateCarHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cars/:id", app.getCarHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cars", app.listCarsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cars/:id", app.deleteCarHandler)
	return router
}
