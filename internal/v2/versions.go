package v2

import (
	"database/sql"
	"os"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type Api struct {
	Orders         OrderV2
	Products       ProductV2
	Collections    CollectionV2
	Variants       VariantV2
	Inventory      InventoryV2
	Fulfillments   FulfillmentV2
	Locations      LocationV2
	Metafields     MetafieldV2
	BulkOperations BulkOperationV2
	Webhooks       WebhookV2
	Services       ServiceV2
}

func V2(db *sql.DB, client *graphify.Client) Api {
	return Api{
		Orders:         OrderV2{DB: db, Client: client},
		Products:       ProductV2{DB: db, Client: client},
		Collections:    CollectionV2{DB: db, Client: client},
		Variants:       VariantV2{DB: db, Client: client},
		Inventory:      InventoryV2{DB: db, Client: client},
		Fulfillments:   FulfillmentV2{DB: db, Client: client},
		Locations:      LocationV2{DB: db, Client: client},
		Metafields:     MetafieldV2{DB: db, Client: client},
		BulkOperations: BulkOperationV2{DB: db, Client: client},
		Webhooks:       WebhookV2{DB: db, Client: client},
		Services:       ServiceV2{DB: db, Client: client},
	}

}

func ShopifyV2(client *graphify.Client) *graphify.Client {
	graphify.NewClient(os.Getenv("STORE_NAME"), graphify.WithToken(os.Getenv("STORE_PASSWORD")), graphify.WithVersion("2023-07"), graphify.WithRetries(5))

	return client
}
