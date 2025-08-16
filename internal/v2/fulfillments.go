package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/graph/model"
	"github.com/pistolricks/cosmetics-api/internal/services"
)

type FulfillmentService interface {
	Create(ctx context.Context, input model.FulfillmentV2Input) error
}

type FulfillmentServiceOp struct {
	Client *services.ClientApi
}

var _ FulfillmentService = &FulfillmentServiceOp{}

type mutationFulfillmentCreateV2 struct {
	FulfillmentCreateV2Result struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"fulfillmentCreateV2(fulfillment: $fulfillment)" json:"fulfillmentCreateV2"`
}

type FulfillmentV2 struct {
	DB     *sql.DB
	Client *services.ClientApi
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
