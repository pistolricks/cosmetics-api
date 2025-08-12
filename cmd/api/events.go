package main

import (
	"net/http"
)

func (app *application) getEventHandler(w http.ResponseWriter, r *http.Request) {

	res := app.chromium.Chrome.GetEventDetails()

	err := app.writeJSON(w, http.StatusOK, envelope{"response": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
