package v2

import (
	"context"
	"database/sql"
	"fmt"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type CollectionV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (v2 CollectionV2) Collections() {

	collections, err := v2.Client.Collection.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, c := range collections {
		fmt.Println(c.Handle)
	}
}
