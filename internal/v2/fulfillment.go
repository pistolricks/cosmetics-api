package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/internal/services"
	"github.com/vinhluan/go-graphql-client"
	"github.com/vinhluan/go-shopify-graphql/model"
)

type FulfillmentService interface {
	Create(ctx context.Context, input model.FulfillmentV2Input) error
}

type FulfillmentServiceOp struct {
	DB     *sql.DB
	Client *services.ClientApi
}

var _ FulfillmentService = &FulfillmentServiceOp{}

type mutationFulfillmentCreateV2 struct {
	FulfillmentCreateV2Result struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"fulfillmentCreateV2(fulfillment: $fulfillment)" json:"fulfillmentCreateV2"`
}

func (s *FulfillmentServiceOp) Create(ctx context.Context, fulfillment model.FulfillmentV2Input) error {
	m := mutationFulfillmentCreateV2{}

	vars := map[string]interface{}{
		"fulfillment": fulfillment,
	}
	err := s.Client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.FulfillmentCreateV2Result.UserErrors) > 0 {
		return fmt.Errorf("UserErrors: %+v", m.FulfillmentCreateV2Result.UserErrors)
	}

	return nil
}

type mutationFulfillmentTrackingInfoUpdate struct {
	FulfillmentTrackingInfoUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"fulfillmentTrackingInfoUpdate(fulfillmentId: $fulfillmentId, trackingInfoInput: $trackingInfoInput, notifyCustomer: $notifyCustomer)" json:"fulfillmentTrackingInfoUpdate"`
}

func (s *FulfillmentServiceOp) FulfillmentTrackingUpdate(ctx context.Context, fulfillmentId string, trackingInfoInput model.FulfillmentTrackingInput, notifyCustomer graphql.Boolean) error {
	m := mutationFulfillmentTrackingInfoUpdate{}

	vars := map[string]interface{}{
		"fulfillmentId":     fulfillmentId,
		"trackingInfoInput": trackingInfoInput,
		"notifyCustomer":    notifyCustomer,
	}
	err := s.Client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.FulfillmentTrackingInfoUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("UserErrors: %+v", m.FulfillmentTrackingInfoUpdateResult.UserErrors)
	}

	return nil
}
