package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
)

type RimanOrder struct {
	OrderDate               string      `json:"orderDate"`
	MainOrdersPK            int         `json:"mainOrdersPK"`
	OrderType               string      `json:"orderType"`
	FinalOrderType          interface{} `json:"finalOrderType"`
	SiteURL                 string      `json:"siteURL"`
	EncOrderNumber          string      `json:"encOrderNumber"`
	CurrencySymbol          string      `json:"currencySymbol"`
	CurrencyCode            string      `json:"currencyCode"`
	PaidStatus              bool        `json:"paidStatus"`
	HasTaxInvoice           bool        `json:"hasTaxInvoice"`
	HasCommercialInvoice    bool        `json:"hasCommercialInvoice"`
	HasCreditNote           bool        `json:"hasCreditNote"`
	IsShippingPending       bool        `json:"isShippingPending"`
	IsPB                    bool        `json:"isPB"`
	IsPA                    bool        `json:"isPA"`
	IsCC                    bool        `json:"isCC"`
	MainFK                  int         `json:"mainFK"`
	MainOrderTypeFK         int         `json:"mainOrderTypeFK"`
	VoucherURL              interface{} `json:"voucherURL"`
	ShipCountry             string      `json:"shipCountry"`
	ShippingStatus          string      `json:"shippingStatus"`
	OrderShippingStatus     string      `json:"orderShippingStatus"`
	OrderTypeValue          interface{} `json:"orderTypeValue"`
	PaidStatusValue         string      `json:"paidStatusValue"`
	Quantity                int         `json:"quantity"`
	Email                   interface{} `json:"email"`
	Phone                   interface{} `json:"phone"`
	ShipFirstName           interface{} `json:"shipFirstName"`
	ShipLastName            interface{} `json:"shipLastName"`
	MarkedPaidDate          string      `json:"markedPaidDate"`
	Total                   float64     `json:"total"`
	ConvTotal               float64     `json:"convTotal"`
	ConvTotalFormat         string      `json:"convTotalFormat"`
	SubTotal                int         `json:"subTotal"`
	ConvSubtotal            int         `json:"convSubtotal"`
	ShipTotal               float64     `json:"shipTotal"`
	ConvShipTotal           float64     `json:"convShipTotal"`
	Taxes                   float64     `json:"taxes"`
	TaxLabel                string      `json:"taxLabel"`
	ProductTax              float64     `json:"productTax"`
	ShippingTax             float64     `json:"shippingTax"`
	AdditionalTaxLabel      string      `json:"additionalTaxLabel"`
	AdditionalTax           interface{} `json:"additionalTax"`
	ConvTaxes               float64     `json:"convTaxes"`
	OrderProcessingFees     interface{} `json:"orderProcessingFees"`
	ConvOrderProcessingFees interface{} `json:"convOrderProcessingFees"`
	Discount                int         `json:"discount"`
	ConvDiscount            int         `json:"convDiscount"`
	RefundAmount            int         `json:"refundAmount"`
	ConvRefund              int         `json:"convRefund"`
	SalesCampaignFK         interface{} `json:"salesCampaignFK"`
	Paidstatusfk            int         `json:"paidstatusfk"`
	DeliveryDate            interface{} `json:"deliveryDate"`
	ShippingDetails         interface{} `json:"shippingDetails"`
	OrderItems              interface{} `json:"orderItems"`
	Payments                interface{} `json:"payments"`
	IsPrepaidOrder          interface{} `json:"isPrepaidOrder"`
	ShowInvoice             bool        `json:"showInvoice"`
	ShowOrderInvoice        bool        `json:"showOrderInvoice"`
	KrGuaranteeNo           string      `json:"krGuaranteeNo"`
	WeChatOrderNumber       interface{} `json:"weChatOrderNumber"`
}

func (app *application) getShipmentHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		OrderId string `json:"order_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	orderId := input.OrderId

	shipment, err := riman.ShipmentHandler(orderId)

	err = app.writeJSON(w, http.StatusOK, envelope{"shipment": shipment, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackingHandler(w http.ResponseWriter, r *http.Request) {
	// https://cart-api.riman.com/api/v1/orders/{rid}/shipment-products

	ctx := r.Context()

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
		Token         string
		OrderId       string
		FulfillmentId string
		ID            string
		data.Filters
	}

	qs := r.URL.Query()
	input.Token = app.readString(qs, "token", "")
	input.OrderId = app.readString(qs, "order_id", "")
	input.FulfillmentId = app.readString(qs, "fulfillment_id", "")
	input.ID = app.readString(qs, "id", "")
	tracking, _ := data.OrderUpdateTracking(input.OrderId, input.Token)

	fId, err := strconv.ParseUint(input.FulfillmentId, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	fulfillment := goshopify.Fulfillment{}

	id, err := strconv.ParseUint(input.ID, 10, 64)

	if tracking[0].TrackingNumber == "" {
		return
	} else {

		fulfillment = goshopify.Fulfillment{
			Id:              fId,
			TrackingUrl:     tracking[0].TrackingLink,
			TrackingNumber:  tracking[0].TrackingNumber,
			TrackingCompany: "Other",
		}
		_, err := client.Order.UpdateFulfillment(ctx, id, fulfillment)

		if err != nil {
			fmt.Println(err.Error())
		}

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tracking": tracking, "fulfillment": fulfillment, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) shippingHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token    string
		UserName string
		data.Filters
	}

	qs := r.URL.Query()
	input.Token = app.readString(qs, "token", "")
	input.UserName = app.readString(qs, "userName", "")

	u := &url.URL{
		Scheme: "https",
		Host:   "cart-api.riman.com",
		Path:   "/api/v1/orders",
	}

	q := u.Query()

	q.Add("mainSiteUrl", input.UserName)
	q.Add("getEnrollerOrders", "")
	q.Add("getCustomerOrders", "")
	q.Add("orderNumber", "")
	q.Add("shipmentNumber", "")
	q.Add("trackingNumber", "")
	q.Add("isRefunded", "")
	q.Add("paidStatus", "true")
	q.Add("orderType", "")
	q.Add("orderLevel", "")
	q.Add("weChatOrderNumber", "")
	q.Add("startDate", "")
	q.Add("endDate", "")
	q.Add("offset", "0")
	q.Add("limit", "40")
	q.Add("orderBy", "-mainOrdersPK")
	q.Add("token", input.Token)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	bodyString := string(resBody)

	err = app.writeJSON(w, http.StatusOK, envelope{"body": bodyString}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
