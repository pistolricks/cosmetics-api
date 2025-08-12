package main

import (
	"context"
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/joho/godotenv"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/shopify"
)

/* ORDER STATUS */
// - open
// - closed
// - cancelled
// - any

/* ORDER FULFILLMENT STATUS */
// - shipped
// - partial
// - unshipped
// - any
// - unfulfilled
// - fulfilled

/* ORDER FINANCIAL STATUS */
// authorized
// pending
// paid
// partially_paid
// refunded
// voided
// partially_refunded
// any
// unpaid

func (app *application) orderUpdateHandler(w http.ResponseWriter, r *http.Request) {
	/*
		shopApp := goshopify.App{
			ApiKey:      app.envars.ShopifyKey,
			ApiSecret:   app.envars.ShopifySecret,
			RedirectUrl: "https://example.com/callback",
			Scope:       "write_orders",
		}

		client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)
		if err != nil {
			fmt.Println(err)
			return
		}

	*/
}

func (app *application) listShopifyOrdersByStatusHandler(w http.ResponseWriter, r *http.Request) {
	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := app.readStringParam("status", r)
	options := struct {
		Status string `url:"status"`
	}{s}

	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listShopifyOrdersByAllStatusValuesHandler(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	s := app.readStringParam("status", r)
	f := app.readStringParam("fulfillment_status", r)

	options := struct {
		Status            string `url:"status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{s, f}

	count, err := client.Order.Count(context.Background(), options)

	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) processShopifyOrder(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Orders []goshopify.Order `json:"orders"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	orders := input.Orders

	count := len(input.Orders)

	rimanStoreName := app.envars.RimanStoreName
	if rimanStoreName == "" {
		fmt.Println("missing riman store name")
		return
	}

	app.background(func() {
		app.ProcessOrders(rimanStoreName, app.page, app.browser, app.cookies, orders)
	})

	currentBrowser := app.browser
	currentPage := app.page
	currentCookies := app.cookies

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": input.Orders, "count": count, "page": currentPage, "browser": currentBrowser, "cookies": currentCookies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) processShopifyOrders(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	options := struct {
		Status            string `url:"status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{"open", "unfulfilled"}

	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	rimanStoreName := app.envars.RimanStoreName // os.Getenv("RIMAN_STORE_NAME")
	if rimanStoreName == "" {
		fmt.Println("missing riman store name")
		return
	}

	loginUrl := app.envars.LoginUrl // os.Getenv("LOGIN_URL")
	if loginUrl == "" {
		fmt.Println("missing login url")
		return
	}

	username := app.envars.Username // os.Getenv("USERNAME")
	if username == "" {
		fmt.Println("missing username")
		return
	}

	password := app.envars.Password // os.Getenv("PASSWORD")
	if password == "" {
		fmt.Println("missing password")
		return
	}

	app.background(func() {
		app.ProcessOrders(rimanStoreName, app.page, app.browser, app.cookies, orders)
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listShopifyOrdersHandler(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	var input struct {
		Status            string
		FinancialStatus   string
		FulfillmentStatus string
		data.Filters
	}

	qs := r.URL.Query()

	fmt.Println(qs)

	input.Status = app.readString(qs, "status", "any")
	input.FinancialStatus = app.readString(qs, "financial_status", "any")
	input.FulfillmentStatus = app.readString(qs, "fulfillment_status", "any")

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id"}

	options := struct {
		Status            string `url:"status"`
		FinancialStatus   string `url:"financial_status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{input.Status, input.FinancialStatus, input.FulfillmentStatus}

	fmt.Println(options)
	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAllShopifyOrders(w http.ResponseWriter, r *http.Request) {
	collection, err := shopify.ListAllOrders()

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": collection}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listShopifyOrders(w http.ResponseWriter, r *http.Request) {
	collection, err := shopify.ListOrders()

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": collection}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listRimanOrders(w http.ResponseWriter, r *http.Request) {

	orderResponse, err := riman.GetOrders(app.envars.Username, app.envars.Token, app.cookies)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	orders := orderResponse.Orders
	count := orderResponse.TotalCount

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
