package shopify

import (
	"github.com/vinhluan/go-shopify-graphql/model"
)

type mutationFulfillmentCreateV2 struct {
	FulfillmentCreateV2Result struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"fulfillmentCreateV2(fulfillment: $fulfillment)" json:"fulfillmentCreateV2"`
}

/* type FulfillmentOrderLineItemsInput struct {
	 FulfillmentOrderID string `json:"fulfillmentOrderId"`
	 FulfillmentOrderLineItems []FulfillmentOrderLineItemInput `json:"fulfillmentOrderLineItems,omitempty,omitempty"`
}*/
