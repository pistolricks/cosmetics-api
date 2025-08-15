package v2

import (
	"database/sql"
	
	graphify "github.com/vinhluan/go-shopify-graphql"
)

type ServiceV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (v2 ServiceV2) Select() {

}
