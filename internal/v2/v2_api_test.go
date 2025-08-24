package v2

import (
	"os"
	"testing"

	r0gql "github.com/r0busta/graphql"
)

func TestV2ConstructorUsesEnv(t *testing.T) {
	os.Setenv("SHOPIFY_TOKEN", "dummy-token")
	os.Setenv("STORE_NAME", "example")
	defer os.Unsetenv("SHOPIFY_TOKEN")
	defer os.Unsetenv("STORE_NAME")

	client := V2(nil, r0gql.Client{})
	gql, _ := client.GraphQLClient()
	if gql == nil {
		t.Fatal("expected non-nil GraphQL client from V2")
	}
}
