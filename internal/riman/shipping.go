package riman

import (
	"database/sql"
	"fmt"

	"resty.dev/v3"
)

type ProductTracking struct {
	PackagePk                 int         `json:"packagePk"`
	ProductPk                 int         `json:"productPk"`
	PackageName               string      `json:"packageName"`
	ProductName               string      `json:"productName"`
	IsPackage                 bool        `json:"isPackage"`
	Quantity                  int         `json:"quantity"`
	Cv                        float64     `json:"cv"`
	Sp                        float64     `json:"sp"`
	Price                     float64     `json:"price"`
	FormattedPrice            string      `json:"formattedPrice"`
	CurrencyCode              string      `json:"currencyCode"`
	ShipmentNumber            string      `json:"shipmentNumber"`
	ShipmentStatus            string      `json:"shipmentStatus"`
	ShippedDate               string      `json:"shippedDate"`
	TrackingNumber            string      `json:"trackingNumber"`
	TrackingLink              string      `json:"trackingLink"`
	VideoOrderPackagingInfoPK interface{} `json:"videoOrderPackagingInfoPK"`
}

type Errors struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

type ShipModel struct {
	DB *sql.DB
}

func (m ShipModel) ShipmentTracker(orderId string, token string) ([]*ProductTracking, error) {

	client := resty.New()
	defer client.Close()

	shipmentUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/orders/%s/shipment-products?token=%s", orderId, token)
	// https://cart-api.riman.com/api/v1/orders/2318156/shipment-products?token=
	res, err := client.R().
		// SetAuthToken(token).
		SetResult(&ProductTracking{}). // or SetResult(LoginResponse{}).
		SetError(&Errors{}).           // or SetError(LoginError{}).
		Post(shipmentUrl)

	if err != nil {
		return nil, err
	}

	return res.Result().([]*ProductTracking), err

}
