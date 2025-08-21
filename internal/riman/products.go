package riman

import (
	"database/sql"
	"fmt"

	"resty.dev/v3"
)

type Products struct{}

/**/
type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) GetProducts() (*[]ProductInformation, error) {

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
		SetResult(&[]ProductInformation{}).
		Get("https://cart-api.riman.com/api/v2/products")

	fmt.Println(err, res)
	products := res.Result().(*[]ProductInformation)

	return products, err
}
