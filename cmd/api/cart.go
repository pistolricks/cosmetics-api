package main

import (
	"fmt"
	"net/http"

	"github.com/pistolricks/cosmetics-api/internal/riman"
	"gopkg.in/guregu/null.v4"
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

func (app *application) updateCartHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token           string      `json:"token"`
		CartKey         string      `json:"cart_key"`
		ConfigFk        null.String `json:"config_fk,omitempty"`
		Discount        float64     `json:"discount,omitempty"`
		ExtraFee        float64     `json:"extra_fee,omitempty"`
		MainCartFk      string      `json:"main_cart_fk"`
		MainCartItemsPk int         `json:"main_cart_items_pk,omitempty"`
		ProductFk       int         `json:"product_fk"`
		Quantity        int         `json:"quantity"`
		SetupForAs      bool        `json:"setup_for_as,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	/*
		product, err := app.riman.Products.GetByFk(productFk)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
	*/

	addToCartPayload := &riman.AddProductPayload{
		ConfigFk:        input.ConfigFk,
		Discount:        input.Discount,
		ExtraFee:        input.ExtraFee,
		MainCartFk:      input.MainCartFk,
		MainCartItemsPk: input.MainCartItemsPk,
		ProductFk:       input.ProductFk,
		Quantity:        input.Quantity,
		SetupForAs:      input.SetupForAs,
	}

	res, err := riman.AddProductToCart(input.Token, input.CartKey, addToCartPayload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cart": res, "error": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
