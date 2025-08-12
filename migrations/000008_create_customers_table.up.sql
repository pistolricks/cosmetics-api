CREATE TABLE IF NOT EXISTS customers
(
    id
    bigserial
    PRIMARY
    KEY,
    created_at
    timestamp
(
    0
) with time zone NOT NULL DEFAULT NOW(),
    first_name text NOT NULL,
    last_name text NOT NULL,
    company_name text,
    email citext UNIQUE NOT NULL,
    phone text NOT NULL,
    contact_method text NOT NULL,
    subscribed bool NOT NULL DEFAULT false,
    location text,
    default_address_id bigint,
    amount_spent float DEFAULT 0,
    order_count bigint DEFAULT 0,
    customer_since text,
    rfm_group text,
    timeline json,
    notes text,
    tags text[],
    last_order_id bigint
    );

CREATE TABLE IF NOT EXISTS customers_addresses
(
    customer_id
    bigint
    NOT
    NULL
    REFERENCES
    customers
    ON
    DELETE
    CASCADE,
    address_id
    bigint
    NOT
    NULL
    REFERENCES
    addresses
    ON
    DELETE
    CASCADE,
    address_type
    text
    NOT
    NULL,
    details
    text,
    PRIMARY
    KEY
(
    customer_id,
    address_id
)
    );