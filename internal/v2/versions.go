package v2

import (
	"database/sql"

	"github.com/pistolricks/cosmetics-api/internal/services"
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
}

func V2(db *sql.DB, cb *services.ClientApi) Api {
	return Api{
		Orders:         OrderV2{DB: db, Client: cb},
		Products:       ProductV2{DB: db, Client: cb},
		Collections:    CollectionV2{DB: db, Client: cb},
		Variants:       VariantV2{DB: db, Client: cb},
		Inventory:      InventoryV2{DB: db, Client: cb},
		Fulfillments:   FulfillmentV2{DB: db, Client: cb},
		Locations:      LocationV2{DB: db, Client: cb},
		Metafields:     MetafieldV2{DB: db, Client: cb},
		BulkOperations: BulkOperationV2{DB: db, Client: cb},
		Webhooks:       WebhookV2{DB: db, Client: cb},
	}

}
