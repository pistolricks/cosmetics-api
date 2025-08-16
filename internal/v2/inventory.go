package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/internal/services"
	"github.com/vinhluan/go-shopify-graphql/model"
)

type InventoryService interface {
	Update(ctx context.Context, id string, input model.InventoryItemUpdateInput) error
	Adjust(ctx context.Context, locationID string, input []model.InventoryAdjustItemInput) error
	ActivateInventory(ctx context.Context, locationID string, id string) error
}

type InventoryV2 struct {
	DB     *sql.DB
	Client *services.ClientApi
}

type mutationInventoryItemUpdate struct {
	InventoryItemUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryItemUpdate(id: $id, input: $input)" json:"inventoryItemUpdate"`
}

type mutationInventoryBulkAdjustQuantityAtLocation struct {
	InventoryBulkAdjustQuantityAtLocationResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryBulkAdjustQuantityAtLocation(locationId: $locationId, inventoryItemAdjustments: $inventoryItemAdjustments)" json:"inventoryBulkAdjustQuantityAtLocation"`
}

type mutationInventoryActivate struct {
	InventoryActivateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryActivate(inventoryItemId: $itemID, locationId: $locationId)" json:"inventoryActivate"`
}

func (s *InventoryV2) Update(ctx context.Context, id string, input model.InventoryItemUpdateInput) error {
	m := mutationInventoryItemUpdate{}
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	err := s.Client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryItemUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryItemUpdateResult.UserErrors)
	}

	return nil
}

func (s *InventoryV2) Adjust(ctx context.Context, locationID string, input []model.InventoryAdjustItemInput) error {
	m := mutationInventoryBulkAdjustQuantityAtLocation{}
	vars := map[string]interface{}{
		"locationId":               locationID,
		"inventoryItemAdjustments": input,
	}
	err := s.Client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors)
	}

	return nil
}

func (s *InventoryV2) ActivateInventory(ctx context.Context, locationID string, id string) error {
	m := mutationInventoryActivate{}
	vars := map[string]interface{}{
		"itemID":     id,
		"locationId": locationID,
	}
	err := s.Client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryActivateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryActivateResult.UserErrors)
	}

	return nil
}
