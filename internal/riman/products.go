package riman

import (
	"fmt"
)

type Products struct{}

/**/

func GetProducts() (*[]RimanProduct, error) {

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetQueryParams(map[string]string{
			"cartType":    "R",
			"countryCode": "US",
			"culture":     "en-US",
			"isCart":      "true",
			"repSiteUrl":  "rmnsocial",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&[]RimanProduct{}).
		Get("https://cart-api.riman.com/api/v2/products")

	fmt.Println(err, res)
	products := res.Result().(*[]RimanProduct)

	return products, err
}
