package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateShopifyOrderID = errors.New("duplicate shopify_order_id")
)

// ShopifyOrder represents a simplified Shopify order record stored locally.
type ShopifyOrder struct {
	ID                int64     `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	ShopifyOrderID    int64     `json:"shopify_order_id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	FinancialStatus   string    `json:"financial_status"`
	FulfillmentStatus string    `json:"fulfillment_status"`
	Currency          string    `json:"currency"`
	TotalPrice        float64   `json:"total_price"`
	Raw               []byte    `json:"raw"`
	Version           int       `json:"-"`
}

// ShopifyOrderModel wraps the DB handle for shopify_orders operations.
type ShopifyOrderModel struct {
	DB *sql.DB
}

// Insert creates a new shopify_orders row and sets ID, CreatedAt, Version on the provided struct.
func (m ShopifyOrderModel) Insert(o *ShopifyOrder) error {
	query := `
        INSERT INTO shopify_orders 
            (shopify_order_id, name, email, financial_status, fulfillment_status, currency, total_price, raw)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, version`

	args := []any{
		o.ShopifyOrderID,
		o.Name,
		o.Email,
		o.FinancialStatus,
		o.FulfillmentStatus,
		o.Currency,
		o.TotalPrice,
		o.Raw,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&o.ID, &o.CreatedAt, &o.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "shopify_orders_shopify_order_id_key"`:
			return ErrDuplicateShopifyOrderID
		default:
			return err
		}
	}

	return nil
}

// GetByID fetches an order by the internal ID.
func (m ShopifyOrderModel) GetByID(id int64) (*ShopifyOrder, error) {
	query := `
        SELECT id, created_at, shopify_order_id, name, email, financial_status, fulfillment_status, currency, total_price, raw, version
        FROM shopify_orders
        WHERE id = $1`

	var o ShopifyOrder

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&o.ID,
		&o.CreatedAt,
		&o.ShopifyOrderID,
		&o.Name,
		&o.Email,
		&o.FinancialStatus,
		&o.FulfillmentStatus,
		&o.Currency,
		&o.TotalPrice,
		&o.Raw,
		&o.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &o, nil
}

// GetByShopifyOrderID fetches an order by the Shopify order ID.
func (m ShopifyOrderModel) GetByShopifyOrderID(shopifyOrderID int64) (*ShopifyOrder, error) {
	query := `
        SELECT id, created_at, shopify_order_id, name, email, financial_status, fulfillment_status, currency, total_price, raw, version
        FROM shopify_orders
        WHERE shopify_order_id = $1`

	var o ShopifyOrder

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, shopifyOrderID).Scan(
		&o.ID,
		&o.CreatedAt,
		&o.ShopifyOrderID,
		&o.Name,
		&o.Email,
		&o.FinancialStatus,
		&o.FulfillmentStatus,
		&o.Currency,
		&o.TotalPrice,
		&o.Raw,
		&o.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &o, nil
}

// Update modifies an existing order using optimistic locking on version.
func (m ShopifyOrderModel) Update(o *ShopifyOrder) error {
	query := `
        UPDATE shopify_orders
        SET name = $1,
            email = $2,
            financial_status = $3,
            fulfillment_status = $4,
            currency = $5,
            total_price = $6,
            raw = $7,
            version = version + 1
        WHERE id = $8 AND version = $9
        RETURNING version`

	args := []any{
		o.Name,
		o.Email,
		o.FinancialStatus,
		o.FulfillmentStatus,
		o.Currency,
		o.TotalPrice,
		o.Raw,
		o.ID,
		o.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&o.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "shopify_orders_shopify_order_id_key"`:
			return ErrDuplicateShopifyOrderID
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete removes an order by ID.
func (m ShopifyOrderModel) Delete(id int64) error {
	query := `
        DELETE FROM shopify_orders
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
