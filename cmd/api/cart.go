package main

import (
	"fmt"
	"net/http"

	"github.com/pistolricks/kbeauty-api/internal/riman"
)

func (app *application) getCartHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		CartKey string `json:"cart_key"`
	}

	fmt.Println(app.session.CartKey)
	fmt.Println(input.CartKey)

	cart, err := riman.GetCart(input.CartKey)

	err = app.writeJSON(w, http.StatusOK, envelope{"cart": cart, "error": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
