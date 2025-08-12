package riman

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

type ShipmentResponse struct {
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

func ShipmentHandler(orderId string) (*ShipmentResponse, error) {

	client := resty.New()
	defer client.Close()

	shipmentUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/orders/%s/shipment-products", orderId)

	token := os.Getenv("TOKEN")

	res, err := client.R().
		SetAuthToken(token).
		SetResult(&ShipmentResponse{}). // or SetResult(LoginResponse{}).
		SetError(&Errors{}).            // or SetError(LoginError{}).
		Get(shipmentUrl)

	if err != nil {
		return nil, err
	}

	return res.Result().(*ShipmentResponse), err

}
