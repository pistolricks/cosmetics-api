package v2

import (
	"context"
	"database/sql"
	"fmt"

	graphify "github.com/vinhluan/go-shopify-graphql"
)

type OrderV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (v2 OrderV2) Orders() {
	
	orders, err := v2.Client.Order.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, c := range orders {
		fmt.Println(c)
	}
}
