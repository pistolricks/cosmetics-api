package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/validator"
)

func (app *application) createRimanSessionHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
		CartKey  string `json:"cartKey"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.UserName)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	client, err := app.riman.Clients.GetByClientUserName(input.UserName)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response, err := app.riman.Session.Login(input.UserName, input.Password, client.Token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	session, err := app.riman.Session.NewRimanSession(client.ID, 24*time.Hour, riman.ScopeAuthentication, response.Jwt, input.CartKey, envelope{"cookies": response.Status, "li_token": response.LiToken, "li_user": response.LiUser, "security_redirect": response.SecurityRedirect})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
