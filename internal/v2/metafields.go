package v2

import (
	"database/sql"
	"fmt"
	"strings"

	"context"

	"github.com/pistolricks/cosmetics-api/graph/model"
	"github.com/pistolricks/cosmetics-api/internal/services"
	log "github.com/sirupsen/logrus"
	graphify "github.com/vinhluan/go-shopify-graphql"
)

type MetafieldService interface {
	ListAllShopMetafields(ctx context.Context) ([]*model.Metafield, error)
	ListShopMetafieldsByNamespace(ctx context.Context, namespace string) ([]*model.Metafield, error)

	GetShopMetafieldByKey(ctx context.Context, namespace, key string) (*model.Metafield, error)

	Delete(ctx context.Context, metafield model.MetafieldDeleteInput) error
	DeleteBulk(ctx context.Context, metafield []model.MetafieldDeleteInput) error
}

type MetafieldServiceOp struct {
	client *services.Client
}

var _ MetafieldService = &MetafieldServiceOp{}

type mutationMetafieldDelete struct {
	MetafieldDeleteResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"metafieldDelete(input: $input)" json:"metafieldDelete"`
}

type MetafieldV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (s *MetafieldServiceOp) ListAllShopMetafields(ctx context.Context) ([]*model.Metafield, error) {
	q := `
		{
			shop{
				metafields{
					edges{
						node{
							createdAt
							description
							id
							key
							legacyResourceId
							namespace
							ownerType
							updatedAt
							value
							type
						}
					}
				}
			}
		}
`

	res := []*model.Metafield{}
	err := s.client.BulkOperations.BulkQuery(ctx, q, &res)
	if err != nil {
		return nil, fmt.Errorf("bulk query: %w", err)
	}

	return res, nil
}

func (s *MetafieldServiceOp) ListShopMetafieldsByNamespace(ctx context.Context, namespace string) ([]*model.Metafield, error) {
	q := `
		{
			shop{
				metafields(namespace: "$namespace"){
					edges{
						node{
							createdAt
							description
							id
							key
							legacyResourceId
							namespace
							ownerType
							updatedAt
							value
							type
						}
					}
				}
			}
		}
`
	q = strings.ReplaceAll(q, "$namespace", namespace)

	res := []*model.Metafield{}
	err := s.client.BulkOperations.BulkQuery(ctx, q, &res)
	if err != nil {
		return nil, fmt.Errorf("bulk query: %w", err)
	}

	return res, nil
}

func (s *MetafieldServiceOp) GetShopMetafieldByKey(ctx context.Context, namespace, key string) (*model.Metafield, error) {
	var q struct {
		Shop struct {
			Metafield model.Metafield `graphql:"metafield(namespace: $namespace, key: $key)"`
		} `graphql:"shop"`
	}
	vars := map[string]interface{}{
		"namespace": namespace,
		"key":       key,
	}

	err := s.client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &q.Shop.Metafield, nil
}

func (s *MetafieldServiceOp) DeleteBulk(ctx context.Context, metafields []model.MetafieldDeleteInput) error {
	for _, m := range metafields {
		err := s.Delete(ctx, m)
		if err != nil {
			log.Warnf("Couldn't delete metafield (%v): %s", m, err)
		}
	}

	return nil
}

func (s *MetafieldServiceOp) Delete(ctx context.Context, metafield model.MetafieldDeleteInput) error {
	m := mutationMetafieldDelete{}

	vars := map[string]interface{}{
		"input": metafield,
	}
	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.MetafieldDeleteResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.MetafieldDeleteResult.UserErrors)
	}

	return nil
}
