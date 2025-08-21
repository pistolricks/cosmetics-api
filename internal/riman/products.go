package riman

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"resty.dev/v3"
)

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

// SaveProducts saves or updates the given products into the riman_products table.
// It stores the entire product payload as JSONB along with the product_pk as the primary key.
func (m ProductModel) SaveProducts(products *[]ProductInformation) (int, error) {
	if products == nil || len(*products) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO riman_products (product_pk, data, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (product_pk)
		DO UPDATE SET data = EXCLUDED.data, updated_at = NOW()`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	saved := 0
	for _, p := range *products {
		b, mErr := json.Marshal(p)
		if mErr != nil {
			err = mErr
			break
		}
		if _, execErr := stmt.ExecContext(ctx, p.ProductPK, string(b)); execErr != nil {
			err = execErr
			break
		}
		saved++
	}

	if err != nil {
		return 0, err
	}

	if cErr := tx.Commit(); cErr != nil {
		return 0, cErr
	}

	return saved, nil
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
