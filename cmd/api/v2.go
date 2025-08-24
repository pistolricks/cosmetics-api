package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	iv2 "github.com/pistolricks/cosmetics-api/internal/v2"
	"github.com/pistolricks/cosmetics-api/internal/validator"
	"github.com/vinhluan/go-graphql-client"
	"github.com/vinhluan/go-shopify-graphql/model"
)

func (app *application) shopifyV2ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure long-running GraphQL requests don't stall the server/client
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	qs := r.URL.Query()
	v := validator.New()

	query := app.readString(qs, "query", "")
	first := app.readInt(qs, "first", 10, v)
	last := app.readInt(qs, "last", 0, v)
	before := app.readString(qs, "before", "")
	after := app.readString(qs, "after", "")
	reverseStr := app.readString(qs, "reverse", "false")
	reverse, _ := strconv.ParseBool(reverseStr)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	svc := &iv2.OrderServiceOp{Client: app.graphify}
	orders, err := svc.List(ctx, iv2.ListOptions{
		Query:   query,
		First:   first,
		Last:    last,
		After:   after,
		Before:  before,
		Reverse: reverse,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := envelope{"orders": orders}
	if err := app.writeJSON(w, http.StatusOK, resp, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// GET /v2/platform/orders
// Query params: query, first, last, before, after, reverse
func (app *application) shopifyV2ListOrdersAfterCursorHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure long-running GraphQL requests don't stall the server/client
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	qs := r.URL.Query()
	v := validator.New()

	query := app.readString(qs, "query", "")
	first := app.readInt(qs, "first", 10, v)
	last := app.readInt(qs, "last", 0, v)
	before := app.readString(qs, "before", "")
	after := app.readString(qs, "after", "")
	reverseStr := app.readString(qs, "reverse", "false")
	reverse, _ := strconv.ParseBool(reverseStr)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	svc := &iv2.OrderServiceOp{Client: app.graphify}
	orders, firstCursor, lastCursor, err := svc.ListAfterCursor(ctx, iv2.ListOptions{
		Query:   query,
		First:   first,
		Last:    last,
		After:   after,
		Before:  before,
		Reverse: reverse,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := envelope{"orders": orders}
	if firstCursor != nil {
		resp["firstCursor"] = *firstCursor
	}
	if lastCursor != nil {
		resp["lastCursor"] = *lastCursor
	}

	if err := app.writeJSON(w, http.StatusOK, resp, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// POST /v2/platform/fulfillments
// Body: { "fulfillment": model.FulfillmentV2Input }
func (app *application) shopifyV2CreateFulfillmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	var input struct {
		Fulfillment model.FulfillmentV2Input `json:"fulfillment"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	svc := &iv2.FulfillmentServiceOp{Client: app.graphify}
	if err := svc.Create(ctx, input.Fulfillment); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_ = app.writeJSON(w, http.StatusCreated, envelope{"status": "created"}, nil)
}

func (app *application) shopifyV2FulfillmentTrackingUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	var input struct {
		FulfillmentId            string                         `json:"fulfillmentId"`
		NotifyCustomer           graphql.Boolean                `json:"notifyCustomer"`
		FulfillmentTrackingInput model.FulfillmentTrackingInput `json:"trackingInfoInput"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	svc := &iv2.FulfillmentServiceOp{Client: app.graphify}
	if err := svc.FulfillmentTrackingUpdate(ctx, input.FulfillmentId, input.FulfillmentTrackingInput, input.NotifyCustomer); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, envelope{"status": "updated"}, nil)
}

// GET /v2/platform/locations/:id
func (app *application) shopifyV2GetLocationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	id := app.readStringParam("id", r)
	if id == "" {
		app.badRequestResponse(w, r, errors.New("invalid id parameter"))
		return
	}

	svc := &iv2.LocationV2{Client: app.graphify}
	loc, err := svc.Get(ctx, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, envelope{"location": loc}, nil)
}

// GET /v2/platform/shop/metafields?namespace=...
func (app *application) shopifyV2ListShopMetafieldsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	namespace := r.URL.Query().Get("namespace")
	svc := &iv2.MetafieldServiceOp{Client: app.graphify}

	if namespace != "" {
		items, err := svc.ListShopMetafieldsByNamespace(ctx, namespace)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		_ = app.writeJSON(w, http.StatusOK, envelope{"metafields": items, "namespace": namespace}, nil)
		return
	}

	items, err := svc.ListAllShopMetafields(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, envelope{"metafields": items}, nil)
}
