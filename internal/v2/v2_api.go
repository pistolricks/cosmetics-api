package v2

import (
	"database/sql"
	"os"

	"github.com/pistolricks/cosmetics-api/internal/services"
	"github.com/r0busta/graphql"
)

type V2_Api struct {
	DB     *sql.DB
	Client *services.ClientApi
}

// V2 is a minimal constructor used by cmd/api to obtain a services.ClientApi instance.
// The graphql.Client parameter is currently unused; the client is created from environment variables.
func V2(db *sql.DB, _ graphql.Client) services.ClientApi {
	_ = db // currently unused in this simple initializer
	c := services.NewClientWithToken(os.Getenv("SHOPIFY_TOKEN"), os.Getenv("STORE_NAME"))
	return *c
}
