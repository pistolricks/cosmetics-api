package main

import (
	"fmt"
	"net/http"

	"context"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/cosmetics-api/internal/shopify"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

func (app *application) RimanApiListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// create a Resty client

	products, err := app.riman.Products.GetProducts()

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

func (app *application) products(client *graphify.Client) {
	// Get products
	products, err := client.Product.List(context.Background(), "")
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, p := range products {
		fmt.Println(p.Title)
	}
}
