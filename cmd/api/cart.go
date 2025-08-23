package main

import (
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/cosmetics-api/internal/riman"
	"gopkg.in/guregu/null.v4"
)

func (app *application) getCartHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		CartKey string `json:"cart_key"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Println(app.riman.Session.CartKey)
	fmt.Println(input.CartKey)

	cart, err := riman.GetCart(input.CartKey)

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

func (app *application) deleteCartProductHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token   string `json:"token"`
		CartKey string `json:"cart_key"`
		Id      string `json:"id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	s, err := riman.DeleteProductFromCart(input.Token, input.CartKey, input.Id)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"is_deleted": s, "item": input.Id, "cart": input.CartKey, "deleted": true}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateShippingAddress(w http.ResponseWriter, r *http.Request) {

	var input struct {
		CartKey string          `json:"cart_key"`
		Email   string          `json:"email"`
		Order   goshopify.Order `json:"order"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.chromium.Chrome.ProcessShipping(app.background, input.Email, app.cookies, input.Order)

	err = app.writeJSON(w, http.StatusOK, envelope{"shipment": true}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
