package services

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/vinhluan/go-graphql-client"
)

var (
	ErrShopifyApi = errors.New("shopify API error")
)

type Services struct {
	Client *ClientApi
}

func (s Services) NewServices(GQL *graphql.GraphQL, accessToken string, apiKey string, apiBasePath string, retries int, timeout time.Duration, transport http.RoundTripper) Services {
	// Initialize your client here if needed. Keeping it minimal to avoid compilation errors.
	return Services{Client: NewClientWithToken(os.Getenv("SHOPIFY_TOKEN"), os.Getenv("STORE_NAME"))}
}
