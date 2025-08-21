package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TrackingInfo struct {
	PackagePk                 int
	ProductPk                 int
	PackageName               string
	ProductName               string
	IsPackage                 bool
	Quantity                  int
	Cv                        float64
	Sp                        float64
	Price                     float64
	FormattedPrice            string
	CurrencyCode              string
	ShipmentNumber            string
	ShipmentStatus            string
	ShippedDate               string
	TrackingNumber            string
	TrackingLink              string
	VideoOrderPackagingInfoPK string
}

var ErrUnauthorized = errors.New("unauthorized")

func OrderUpdateTracking(orderId string, token string) ([]TrackingInfo, error) {
	path := fmt.Sprintf("/api/v1/orders/%s/shipment-products", orderId)

	u := &url.URL{
		Scheme: "https",
		Host:   "cart-api.riman.com",
		Path:   path,
	}

	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(u.String())
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: received non-200 status code: %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %w", err)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	var trackingInfo []TrackingInfo
	if err := json.Unmarshal(resBody, &trackingInfo); err != nil {
		return nil, fmt.Errorf("client: could not unmarshal response body: %w", err)
	}

	fmt.Println(trackingInfo)

	return trackingInfo, nil
}

func OrderUpdateFulfillment() {

}
