package main

import (
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/kbeauty-api/internal/shopify"
	"github.com/pistolricks/kbeauty-api/internal/vendors"
)

func (app *application) RimanApiListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// create a Resty client

	products, err := vendors.GetProducts()

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) ShopifyApiListProductsHandler(w http.ResponseWriter, r *http.Request) {

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
