package shopify

import (
	"context"
	"fmt"

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

func Create(ctx context.Context, fulfillment model.FulfillmentV2Input) error {
	m := mutationFulfillmentCreateV2{}

	vars := map[string]interface{}{
		"fulfillment": fulfillment,
	}
	err := Client().Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.FulfillmentCreateV2Result.UserErrors) > 0 {
		return fmt.Errorf("UserErrors: %+v", m.FulfillmentCreateV2Result.UserErrors)
	}

	return nil
}
