package main

import (
	"errors"
	"gearboxd/internal/data"
	"gearboxd/internal/validator"
	"net/http"

	"github.com/shopspring/decimal"
)

func (app *application) createCarLogHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CarID   int64            `json:"car_id"`
		Rating  *decimal.Decimal `json:"rating"`
		Status  string           `json:"status"`
		Comment *string          `json:"comment"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var cl data.CarLog
	user := app.contextGetUser(r)
	if user == nil {
		panic("No user found in context for createCarLogHandler")
	}

	cl.UserID = user.ID
	cl.CarID = input.CarID

	if input.Rating != nil {
		cl.Rating = input.Rating
	}

	cl.Status = input.Status

	if input.Comment != nil {
		cl.Comment = input.Comment
	}

	v := validator.New()
	data.ValidateCarLog(v, &cl)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Cars.Get(cl.CarID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.CarLogs.Insert(&cl)
	if err != nil {
		if errors.Is(err, data.ErrCarLogAlreadyExists) {
			app.resourceAlreadyExistsResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"car_log": cl}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCarLogHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CarID int64 `json:"car_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.CarID > 0, "car_id", "cannot be 0 or less")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Cars.Get(input.CarID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)
	if user == nil {
		panic("No user found in context for createCarLogHandler")
	}
	err = app.models.CarLogs.Delete(user.ID, input.CarID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "car log successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCarLogHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.entityNotFoundResponse(w, r)
		return
	}

	carLog, err := app.models.CarLogs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.entityNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	user := app.contextGetUser(r)
	if user.ID != carLog.UserID {
		app.notPermittedResponse(w, r)
	}

	var input struct {
		Rating  *decimal.Decimal `json:"rating"`
		Status  *string          `json:"status"`
		Comment *string          `json:"comment"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if input.Rating != nil {
		carLog.Rating = input.Rating
	}

	if input.Status != nil {
		carLog.Status = *input.Status
	}

	if input.Comment != nil {
		carLog.Comment = input.Comment
	}

	v := validator.New()
	data.ValidateCarLog(v, carLog)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.CarLogs.Update(carLog)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"car_log": carLog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
