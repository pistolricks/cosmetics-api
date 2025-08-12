package riman

import (
	"fmt"
	"net/http"
	"time"

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

func CookieStatus(s proto.NetworkCookieSameSite) (int, error) {
	switch s {
	case "Strict":
		return 3, nil
	case "Lax":
		return 2, nil
	case "None":
		return 4, nil
	// Return a zero value for Status and an error for invalid input.
	default:
		return 1, fmt.Errorf("unknown status: %q", s)
	}
}

func restyCookies(cookies []*proto.NetworkCookie) []*http.Cookie {

	var updatedCookies []*http.Cookie

	for _, cookie := range cookies {

		status, err := CookieStatus(cookie.SameSite)
		if err != nil {
			fmt.Println(err)
		}

		var epochSeconds int64 = int64(cookie.Expires)

		t := time.Unix(epochSeconds, 0)

		updatedCookies = append(updatedCookies, &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
			SameSite: http.SameSite(status),
			Expires:  t,
		})
	}

	return updatedCookies
}

func GetOrders(username string, token string, cookies []*proto.NetworkCookie) (*OrderResponse, error) {
	client := resty.New()
	defer client.Close()

	mainSiteUrl := username
	updatedCookies := restyCookies(cookies)

	url := fmt.Sprintf("https://cart-api.riman.com/api/v1/orders")
	// url := fmt.Sprintf("https://cart-api.riman.com/api/v1/orders?mainSiteUrl=%s&memberID=&getEnrollerOrders=&getCustomerOrders=&orderNumber=&shipmentNumber=&trackingNumber=&isRefunded=&paidStatus=true&orderType=&orderLevel=&weChatOrderNumber=&startDate=&endDate=&offset=0&limit=10&orderBy=-mainOrdersPK", mainSiteUrl)
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
