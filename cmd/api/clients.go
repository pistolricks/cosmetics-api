package main

import (
	"net/http"
)

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
