package v2

import (
	"database/sql"
	"os"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type Api struct {
	Service     ServiceV2
	Orders      OrderV2
	Products    ProductV2
	Collections CollectionV2
}

func V2(db *sql.DB, client *graphify.Client) Api {
	return Api{
		Orders:   OrderV2{DB: db, Client: client},
		Products: ProductV2{DB: db, Client: client},
	}

}

func ShopifyV2() *graphify.Client {
	client := graphify.NewClient(os.Getenv("STORE_NAME"),
		graphify.WithToken(os.Getenv("STORE_PASSWORD")),
		graphify.WithVersion("2023-07"),
		graphify.WithRetries(5))

	return client
}
