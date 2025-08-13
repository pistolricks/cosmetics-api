package main

import "net/http"

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.riman.Session.Login(input.UserName, input.Password)

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
