CREATE TABLE IF NOT EXISTS sessions (
    hash bytea PRIMARY KEY,
    client_id bigint NOT NULL REFERENCES clients ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL,
    cart_key text NOT NULL,
    data json NOT NULL
);