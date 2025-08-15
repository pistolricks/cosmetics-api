package main

import (
	"fmt"
	"net/http"
)

func (app *application) getCartHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		CartKey string `json:"cart_key"`
		Token   string `json:"token"`
	}

	fmt.Println(app.riman.Session.CartKey)
	fmt.Println(input.CartKey)

	cart, err := app.riman.Clients.Patch(input.CartKey, input.Token)

	err = app.writeJSON(w, http.StatusOK, envelope{"cart": cart, "error": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
