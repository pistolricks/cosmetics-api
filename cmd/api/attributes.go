package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func (app *application) findCookieValue() *string {
	for i := range app.cookies {
		if app.cookies[i].Name == "token" {
			app.envars.Token = app.cookies[i].Value
			fmt.Println("app.envars.Token")
			fmt.Println(app.envars.Token)
			return &app.cookies[i].Value
		}
	}
	// Return nil if no product is found
	return nil
}

func (app *application) updateOrderFields(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID      string `json:"id"`
		OrderId string `json:"order_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	metaRad := goshopify.Metafield{Id: 29021818519600, Value: app.client.Username}

	metaRid := goshopify.Metafield{Id: 29021818519600, Value: input.OrderId}

	updatedRad, err := app.shopify.Orders.Config.Client.Order.UpdateMetafield(context.Background(), id, metaRad)
	if err != nil {
		fmt.Println(err)
		return
	}

	updatedRid, err := app.shopify.Orders.Config.Client.Order.UpdateMetafield(context.Background(), id, metaRid)
	if err != nil {
		fmt.Println(err)
		return
	}

	order, err := app.shopify.Orders.Config.Client.Order.Get(context.Background(), id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	updatedOrder := goshopify.Order{Id: order.Id, Note: input.OrderId}
	fmt.Println("TEST UPDATED ORDER")
	fmt.Println(updatedOrder)

	updated, err := app.shopify.Orders.Config.Client.Order.Update(context.Background(), updatedOrder)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": updated, "rad": updatedRad, "rid": updatedRid}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
