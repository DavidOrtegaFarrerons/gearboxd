package main

import (
	"errors"
	"fmt"
	"gearboxd/internal/data"
	"gearboxd/internal/validator"
	"net/http"

	"github.com/shopspring/decimal"
)

func (app *application) createCarHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Make        string  `json:"make"`
		Model       string  `json:"model"`
		Year        int     `json:"year"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		Gearbox     string  `json:"gearbox"`
		Drivetrain  string  `json:"drivetrain"`
		Horsepower  int     `json:"horsepower"`
		Fuel        string  `json:"fuel"`
		PriceNew    float64 `json:"price_new"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	car := &data.Car{
		Make:        input.Make,
		Model:       input.Model,
		Year:        input.Year,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Gearbox:     input.Gearbox,
		Drivetrain:  input.Drivetrain,
		Horsepower:  input.Horsepower,
		Fuel:        input.Fuel,
		PriceNew:    decimal.NewFromFloat(input.PriceNew),
	}

	v := validator.New()
	if data.ValidateCar(v, car); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cars.Insert(car)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cars/%d", car.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"car": car}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil || id < 1 {
		app.entityNotFoundResponse(w, r)
		return
	}

	car, err := app.models.Cars.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"car": car}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil || id < 1 {
		app.entityNotFoundResponse(w, r)
		return
	}

	err = app.models.Cars.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "car successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
