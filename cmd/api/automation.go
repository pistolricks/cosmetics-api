package main

import (
	"net/http"
	"os"
)

func (app *application) inputShippingHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) inputBillingHandler(w http.ResponseWriter, r *http.Request) {

	isFinished := app.chromium.Chrome.InsertBillingInfo(app.addressClient, os.Getenv("ACCOUNT_EMAIL"))

	err := app.writeJSON(w, http.StatusOK, envelope{"isFinished": isFinished}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
