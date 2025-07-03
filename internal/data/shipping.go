package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TrackingInfo struct {
	packagePk                 int
	productPk                 int
	packageName               string
	productName               string
	isPackage                 bool
	quantity                  int
	cv                        float64
	sp                        float64
	price                     float64
	formattedPrice            string
	currencyCode              string
	shipmentNumber            string
	shipmentStatus            string
	shippedDate               string
	trackingNumber            string
	trackingLink              string
	videoOrderPackagingInfoPK string
}

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
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	fmt.Println(u.String())
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

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

	return trackingInfo, nil

}

func OrderUpdateFulfillment() {

}
