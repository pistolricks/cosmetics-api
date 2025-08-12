package main

import (
	"fmt"
	"net/http"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func (app *application) updateOrderMetaField(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var input struct {
		ID        string `json:"id"`
		AccountId string `json:"account_id"`
		OrderId   string `json:"order_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
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

	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateAccountId := goshopify.Metafield{Id: 29021818519600, Value: input.AccountId}
	updatedAccountId, err := client.Order.UpdateMetafield(ctx, id, updateAccountId)
	if err != nil {
		fmt.Println(err)
		return
	}

	updateOrderId := goshopify.Metafield{Id: 29021818585136, Value: input.OrderId}
	updatedOrderId, err := client.Order.UpdateMetafield(ctx, id, updateOrderId)
	if err != nil {
		fmt.Println(err)
		return
	}

	order, err := client.Order.Get(ctx, id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	updatedOrder := goshopify.Order{Id: order.Id, Note: input.OrderId}
	fmt.Println("TEST UPDATED ORDER")
	fmt.Println(updatedOrder)

	updated, err := client.Order.Update(ctx, updatedOrder)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": updated, "account_id": updatedAccountId, "order_id": updatedOrderId}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
