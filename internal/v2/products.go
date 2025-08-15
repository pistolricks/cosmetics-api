package v2

import (
	"context"
	"database/sql"
	"fmt"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type ProductV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (v2 ProductV2) Collections() {

	collections, err := v2.Client.Product.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, c := range collections {
		fmt.Println(c.Handle)
	}
}
