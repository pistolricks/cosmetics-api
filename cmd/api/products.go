package main

import (
	"fmt"
	"net/http"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/cosmetics-api/internal/shopify"
)

func (app *application) rimanApiListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// create a Resty client
	var input struct {
		Token string
		Shop  string
	}

	qs := r.URL.Query()
	input.Token = app.readString(qs, "token", "")

	if input.Token == "" {
		input.Token = r.Header.Get("Authorization")
		if strings.HasPrefix(input.Token, "Bearer ") {
			input.Token = strings.TrimPrefix(input.Token, "Bearer ")
		}
	}

	token := input.Token
	if token == "" || strings.EqualFold(token, "null") {
		token = "TEST"
	}

	input.Shop = app.readString(qs, "shop", "")

	shop := input.Shop
	if shop == "" || strings.EqualFold(shop, "null") {
		shop = "WeKBeauty"
	}

	products, h, err := app.riman.Products.GetProducts(token, shop)

	_ = h

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "errors": err}, nil)
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
