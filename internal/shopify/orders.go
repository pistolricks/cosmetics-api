package shopify

import (
	"context"
	"fmt"

	"github.com/vinhluan/go-shopify-graphql"
	"github.com/vinhluan/go-shopify-graphql/model"
)

func ListAllOrders() ([]*model.Order, error) {

	// Get all collections
	collections, err := Client().Order.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, c := range collections {
		fmt.Println(c.ID)
	}

	return collections, err
}

func ListOrders() ([]*model.Order, error) {

	// Define options to fetch up to 250 orders, a common limit for Shopify API pagination.
	options := shopify.ListOptions{First: 1}

	orders, err := Client().Order.List(context.Background(), options)
	if err != nil {
		// Return the error to the caller for proper handling instead of panicking.
		return nil, err
	}

	return orders, nil

}
