CREATE TABLE IF NOT EXISTS riman_products
(
    product_pk integer PRIMARY KEY,
    data       jsonb                        NOT NULL,
    updated_at timestamp(0) with time zone  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS riman_products_updated_at_idx ON riman_products (updated_at);
