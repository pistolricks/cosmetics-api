CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS shopify_orders
(
    id                 bigserial PRIMARY KEY,
    created_at         timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    shopify_order_id   bigint UNIQUE               NOT NULL,
    name               text,
    email              citext,
    financial_status   text,
    fulfillment_status text,
    currency           text,
    total_price        float,
    raw                jsonb,
    version            integer                     NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS shopify_orders_email_idx ON shopify_orders (email);
