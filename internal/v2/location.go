package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/internal/services"
	"github.com/vinhluan/go-shopify-graphql/model"
)

//go:generate mockgen -destination=./mock/location_service.go -package=mock . LocationService
type LocationService interface {
	Get(ctx context.Context, id string) (*model.Location, error)
}

type LocationV2 struct {
	DB     *sql.DB
	Client *services.ClientApi
}

var _ LocationService = &LocationV2{}

func (s *LocationV2) Get(ctx context.Context, id string) (*model.Location, error) {
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
	err := s.Client.QueryString(ctx, q, vars, &out)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out.Location, nil
}
