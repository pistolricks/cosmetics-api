package riman

import (
	"context"
	"database/sql"
	"net/http"

	"errors"
	"time"

	"resty.dev/v3"
)

/**/
type ProductModel struct {
	DB *sql.DB
}

type ResponseData = map[string]any
type ResultsResponse map[string]any

func (m ProductModel) GetProducts(token string, shop string) (*[]ProductInformation, *http.Response, error) {

	client := resty.New()
	defer func(client *resty.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	res, err := client.R().
		SetAuthToken(token).
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
		SetHeader("Referer", "https://www.riman.com/").
		SetQueryParams(map[string]string{
			"cartType":    "W",
			"countryCode": "US",
			"culture":     "en-US",
			"isCart":      "true",
			"repSiteUrl":  shop,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&[]ProductInformation{}).
		Get("https://cart-api.riman.com/api/v2/products")

	if err != nil {
		return nil, nil, err
	}

	return res.Result().(*[]ProductInformation), res.RawResponse, nil
}

func (m ProductModel) GetByFk(productPk int64) (*ProductInformation, error) {
	query := `
        SELECT product_pk
        FROM riman_products
        WHERE product_pk = $1`

	var pInformation ProductInformation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, productPk).Scan(
		&pInformation.ProductPK,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &pInformation, nil
}
