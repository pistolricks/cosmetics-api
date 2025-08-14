package riman

import (
	"database/sql"
	"fmt"

	"github.com/go-rod/rod/lib/proto"
	"resty.dev/v3"
)

// https://cart-api.riman.com/api/v1/orders
// ?mainSiteUrl=2043124962
// &getEnrollerOrders=
// &getCustomerOrders=
// &orderNumber=
// &shipmentNumber=
// &trackingNumber=
// &isRefunded=
// &paidStatus=true
// &orderType=
// &orderLevel=
// &weChatOrderNumber=
// &startDate=
// &endDate=
// &offset=0
// &limit=20
// &orderBy=-mainOrdersPK

type OrderResponse struct {
	TotalCount int     `json:"totalCount"`
	Orders     []Order `json:"orders"`
}

type Status int

type OrderModel struct {
	DB *sql.DB
}

func (m OrderModel) GetOrders(username string, token string, cookies []*proto.NetworkCookie) (*OrderResponse, error) {
	client := resty.New()
	defer client.Close()

	mainSiteUrl := username
	updatedCookies := restyCookies(cookies)

	url := fmt.Sprintf("https://cart-api.riman.com/api/v1/orders")
	// url := fmt.Sprintf("
	//https://cart-api.riman.com/api/v1/orders?mainSiteUrl=%s&memberID=&getEnrollerOrders=&getCustomerOrders=&orderNumber=&shipmentNumber=&trackingNumber=&isRefunded=&paidStatus=true&orderType=&orderLevel=&weChatOrderNumber=&startDate=&endDate=&offset=0&limit=10&orderBy=-mainOrdersPK", mainSiteUrl)
	res, err := client.R().
		SetAuthToken(token).
		SetCookies(updatedCookies).
		SetQueryParams(map[string]string{
			"mainSiteUrl":       mainSiteUrl,
			"getEnrollerOrders": "",
			"getCustomerOrders": "",
			"orderNumber":       "",
			"shipmentNumber":    "",
			"trackingNumber":    "",
			"isRefunded":        "",
			"paidStatus":        "true",
			"orderType":         "",
			"orderLevel":        "",
			"weChatOrderNumber": "",
			"startDate":         "",
			"endDate":           "",
			"offset":            "0",
			"limit":             "20",
			"orderBy":           "-mainOrdersPK",
		}).
		SetResult(&OrderResponse{}).
		SetError(&Errors{}).
		Get(url)

	fmt.Println(err, res)
	orderResponse := res.Result().(*OrderResponse)

	fmt.Println(orderResponse.Orders)
	fmt.Println(url)

	return orderResponse, err
}

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
	TotalProductTax         float64     `json:"totalProductTax"`
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
	MemberID                interface{} `json:"memberID"`
}
