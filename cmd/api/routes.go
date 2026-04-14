package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/cars", app.requirePermission("cars:write", app.createCarHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/cars/:id", app.requirePermission("cars:write", app.updateCarHandler))
	router.HandlerFunc(http.MethodGet, "/v1/cars/:id", app.requirePermission("cars:read", app.getCarHandler))
	router.HandlerFunc(http.MethodGet, "/v1/cars", app.requirePermission("cars:read", app.listCarsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/cars/:id", app.requirePermission("cars:write", app.deleteCarHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.enableCORS(app.rateLimit(app.authenticate(router)))
}
