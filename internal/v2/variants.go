package v2

import (
	"database/sql"
	"fmt"

	"context"

	"github.com/pistolricks/cosmetics-api/internal/services"
	"github.com/vinhluan/go-shopify-graphql/model"
)

type VariantService interface {
	Update(ctx context.Context, variant model.ProductVariantInput) error
}

type mutationProductVariantUpdate struct {
	ProductVariantUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productVariantUpdate(input: $input)" json:"productVariantUpdate"`
}

type VariantV2 struct {
	DB     *sql.DB
	Client *services.ClientApi
}

func (s *VariantV2) Update(ctx context.Context, variant model.ProductVariantInput) error {
	m := mutationProductVariantUpdate{}

	vars := map[string]interface{}{
		"input": variant,
	}
	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductVariantUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantUpdateResult.UserErrors)
	}

	return nil
}
