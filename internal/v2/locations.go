package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/graph/model"
	"github.com/pistolricks/cosmetics-api/internal/services"
	graphify "github.com/vinhluan/go-shopify-graphql"
)

type LocationService interface {
	Get(ctx context.Context, id string) (*model.Location, error)
}

type LocationServiceOp struct {
	client *services.Client
}

var _ LocationService = &LocationServiceOp{}

type LocationV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (s *LocationServiceOp) Get(ctx context.Context, id string) (*model.Location, error) {
	q := `query location($id: ID!) {
		location(id: $id){
			id
			name
		}
	}`

	vars := map[string]interface{}{
		"id": id,
	}

	var out struct {
		*model.Location `json:"location"`
	}
	err := s.client.QueryString(ctx, q, vars, &out)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out.Location, nil
}
