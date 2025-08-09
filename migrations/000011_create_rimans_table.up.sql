CREATE TABLE IF NOT EXISTS rimans (
    hash bytea PRIMARY KEY,
    client_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL,
    cart_key text NOT NULL,
    data json NOT NULL
);