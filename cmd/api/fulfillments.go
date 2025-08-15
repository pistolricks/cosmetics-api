package main

import (
	"context"
	"net/http"
	"strconv"
)

func (app *application) updateFulfillmentHandler(w http.ResponseWriter, r *http.Request) {
	// https://cart-api.riman.com/api/v1/orders/{rid}/shipment-products

	var input struct {
		FulfillmentID  string `json:"fulfillment_id"`
		TrackingNumber string `json:"tracking_number"`
		TrackingLink   string `json:"tracking_link"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := strconv.ParseUint(input.FulfillmentID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fulfillment, err := app.shopify.Fulfillments.Config.Client.FulfillmentService.Get(context.Background(), id, "")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"fulfillment_id":  input.FulfillmentID,
		"tracking_number": input.TrackingNumber,
		"tracking_link":   input.TrackingLink,
		"fulfillment":     fulfillment, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) importAndUpdateTrackingHandler(w http.ResponseWriter, r *http.Request) {
	// https://cart-api.riman.com/api/v1/orders/{rid}/shipment-products

	var input struct {
		RimanID        string `json:"riman_id"`
		FulfillmentID  string `json:"fulfillment_id"`
		TrackingNumber string `json:"tracking_number"`
		TrackingLink   string `json:"tracking_link"`
		Token          string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tracking, _ := app.riman.Shipping.ShipmentTracker(input.RimanID, input.Token)

	err = app.writeJSON(w, http.StatusOK, envelope{"tracking": tracking, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
