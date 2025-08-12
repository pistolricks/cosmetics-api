CREATE TABLE IF NOT EXISTS addresses
(
    id
    bigserial
    PRIMARY
    KEY,
    first_name
    text
    NOT
    NULL,
    last_name
    text
    NOT
    NULL,
    company
    text,
    address1
    text
    NOT
    NULL,
    address2
    text,
    city
    text,
    province
    text,
    province_code
    text,
    zip
    text,
    country
    text,
    country_code
    text,
    latitude
    float,
    longitude
    float,
    name
    text,
    phone
    text,
    email
    text
);

CREATE TABLE IF NOT EXISTS clients
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
    middle_name text,
    last_name text NOT NULL,
    suffix text,
    email citext UNIQUE NOT NULL,
    mobile text,
    username text NOT NULL,
    riman_user_id bigint NOT NULL,
    status bool NOT NULL,
    organization_type text,
    signup_date text NOT NULL,
    anniversary_date text NOT NULL,
    account_type text NOT NULL,
    sponsor_username text,
    member_id text NOT NULL,
    rank text,
    enrollment_date text,
    personal_orders_volume float,
    personal_clients_volume float,
    total_personal_volume float,
    current_month_sp float,
    current_month_bp float,
    last_order_date text,
    last_order_id bigint,
    last_order_sp float,
    last_order_bp float,
    lifetime_spend float,
    most_recent_12_month_spend float,
    data json,
    password_hash bytea NOT NULL
    );

CREATE TABLE IF NOT EXISTS clients_addresses
(
    client_id
    bigint
    NOT
    NULL
    REFERENCES
    clients
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
    client_id,
    address_id
)
    );