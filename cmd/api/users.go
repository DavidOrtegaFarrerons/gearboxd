package main

import (
	"errors"
	"gearboxd/internal/data"
	"gearboxd/internal/validator"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := &data.User{
		Email:     input.Email,
		Username:  input.Username,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateUser(user, v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddErrorKey("username", "already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddErrorKey("email", "already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
