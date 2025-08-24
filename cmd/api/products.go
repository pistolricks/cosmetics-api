package main

import (
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/cosmetics-api/internal/shopify"
	v2 "github.com/pistolricks/cosmetics-api/internal/v2"
)

func (app *application) rimanApiListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// create a Resty client

	products, err := app.riman.Products.GetProducts()

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// SaveRimanProductsHandler fetches products from the Riman API and saves them into Postgres.
func (app *application) saveRimanProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.riman.Products.GetProducts()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	saved, err := app.riman.Products.SaveProducts(products)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"saved": saved}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) shopifyApiListProductsHandler(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_products",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	products, count, err := shopify.GetProducts(client)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "count": count, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) productsHandler(w http.ResponseWriter, r *http.Request) {
	// Get products

	products, err := v2.Products()

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
