package main

import (
	"net/http"
	"strconv"

	"github.com/pistolricks/cosmetics-api/internal/riman"
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

	fulfillment, err := app.shopify.Fulfillments.UpdateFulfillment(id, input.TrackingNumber, input.TrackingLink)
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
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tracking, _ := app.riman.Shipping.ShipmentTracker(input.RimanID, app.session.Plaintext)

	var trackData riman.ProductTracking

	if len(tracking) > 0 {
		trackData = *tracking[0] // this copies all fields, like JS {...tracking[0]}

		id, err := strconv.ParseUint(input.FulfillmentID, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		fulfillment, err := app.shopify.Fulfillments.UpdateFulfillment(id, trackData.TrackingNumber, trackData.TrackingLink)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"tracking": tracking, "fulfillment": fulfillment, "errors": err}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}

	}
}
