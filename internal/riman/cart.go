package riman

import (
	"fmt"

	"resty.dev/v3"
)

type CartErrors struct {
	Error string `json:"error"`
}

type Body map[string]any

func (m ClientModel) Patch(cartKey string, token string) (*Cart, error) {

	cartUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s", "")

	fmt.Println(cartUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetAuthToken("").
		SetHeader("Accept", "application/json").
		SetBody(&Body{
			"CountryCode": "US",
			"Culture":     "en-US",
		}).
		SetResult(&Cart{}).
		SetError(&CartErrors{}).
		Patch(cartUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*Cart))

	return res.Result().(*Cart), err
}
