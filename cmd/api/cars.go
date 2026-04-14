package main

import (
	"errors"
	"fmt"
	"gearboxd/internal/data"
	"gearboxd/internal/validator"
	"net/http"
	"strconv"

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
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"car": car}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCarHandler(w http.ResponseWriter, r *http.Request) {
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
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(car.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Make        *string          `json:"make"`
		Model       *string          `json:"model"`
		Year        *int             `json:"year"`
		Description *string          `json:"description"`
		ImageURL    *string          `json:"image_url"`
		Gearbox     *string          `json:"gearbox"`
		Drivetrain  *string          `json:"drivetrain"`
		Horsepower  *int             `json:"horsepower"`
		Fuel        *string          `json:"fuel"`
		PriceNew    *decimal.Decimal `json:"price_new"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Make != nil {
		car.Make = *input.Make
	}
	if input.Model != nil {
		car.Model = *input.Model
	}
	if input.Year != nil {
		car.Year = *input.Year
	}
	if input.Description != nil {
		car.Description = *input.Description
	}
	if input.ImageURL != nil {
		car.ImageURL = *input.ImageURL
	}
	if input.Gearbox != nil {
		car.Gearbox = *input.Gearbox
	}
	if input.Drivetrain != nil {
		car.Drivetrain = *input.Drivetrain
	}
	if input.Horsepower != nil {
		car.Horsepower = *input.Horsepower
	}
	if input.Fuel != nil {
		car.Fuel = *input.Fuel
	}
	if input.PriceNew != nil {
		car.PriceNew = *input.PriceNew
	}

	v := validator.New()
	if data.ValidateCar(v, car); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cars.Update(car)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

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

func (app *application) listCarsHandler(w http.ResponseWriter, r *http.Request) {
	//Filter by: make, year, gearbox, drivetrain, fuel horsepower_min / horsepower_max price_min / price_max
	//Sort by: make, year, horsepower, price
	var cf data.CarFilters

	qs := r.URL.Query()
	v := validator.New()

	cf.Make = app.readQueryString(qs, "make", "")
	cf.Model = app.readQueryString(qs, "model", "")
	cf.Year = app.readQueryInt(qs, "year", 0, v)
	cf.Gearbox = app.readQueryString(qs, "gearbox", "")
	cf.Drivetrain = app.readQueryString(qs, "drivetrain", "")
	cf.Fuel = app.readQueryString(qs, "fuel", "")
	cf.HorsepowerMin = app.readQueryInt(qs, "horsepower_min", 0, v)
	cf.HorsepowerMax = app.readQueryInt(qs, "horsepower_max", 0, v)
	cf.PriceMin = app.readQueryDecimal(qs, "price_min", decimal.Zero, v)
	cf.PriceMax = app.readQueryDecimal(qs, "price_max", decimal.Zero, v)

	cf.Filters.Page = app.readQueryInt(qs, "page", 1, v)
	cf.Filters.PageSize = app.readQueryInt(qs, "page_size", 20, v)
	cf.Filters.Sort = app.readQueryString(qs, "sort", "make")
	cf.Filters.SortSafelist = []string{"make", "-make", "year", "-year", "horsepower", "-horsepower", "price", "-price"}

	if data.ValidateFilters(v, cf.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cars, metadata, err := app.models.Cars.GetAll(&cf)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cars": cars, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
