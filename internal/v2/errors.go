package v2

import (
	"database/sql"
	"strings"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type ErrorsV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func IsConnectionError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "broken pipe"))
}
