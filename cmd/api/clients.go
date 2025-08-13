package main

import (
	"net/http"
)

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

func (app *application) findCartKeyValue() *string {
	for i := range app.cookies {
		if app.cookies[i].Name == "cartKey" {
			return &app.cookies[i].Value
		}
	}
	// Return nil if no product is found
	return nil
}

func (app *application) listClientsHandler(w http.ResponseWriter, r *http.Request) {

	clients, metadata, err := app.riman.Clients.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"clients": clients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
