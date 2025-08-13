package main

import (
	"net/http"
)

func (app *application) inputBillingHandler(w http.ResponseWriter, r *http.Request) {

	isFinished := app.chromium.Chrome.Insert()

	err := app.writeJSON(w, http.StatusOK, envelope{"isFinished": isFinished}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
