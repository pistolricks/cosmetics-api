CREATE TABLE IF NOT EXISTS teams
(
    id
    bigint
    PRIMARY
    KEY,
    data
    json,
    clients
    json
)