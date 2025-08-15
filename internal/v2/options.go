package v2

import (
	"fmt"
	"net/http"
	"time"

	graphify "github.com/vinhluan/go-graphql-client"
)

type Option func(shopClient *Client)

func WithGraphQLClient(gql graphify.GraphQL) Option {
	return func(c *Client) {
		c.gql = gql
	}
}

// WithVersion optionally sets the API version if the passed string is valid.
func WithVersion(apiVersion string) Option {
	return func(c *Client) {
		if apiVersion != "" {
			c.apiBasePath = fmt.Sprintf("%s/%s", defaultAPIBasePath, apiVersion)
		}
	}
}

// WithToken optionally sets access token.
func WithToken(token string) Option {
	return func(c *Client) {
		c.accessToken = token
	}
}

// WithPrivateAppAuth optionally sets private app credentials (API key and access token).
func WithPrivateAppAuth(apiKey string, accessToken string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
		c.accessToken = accessToken
	}
}

// WithRetries optionally sets maximum retry count for an API call.
func WithRetries(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

// WithTimeout optionally sets timeout for each HTTP requests made.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithTransport optionally sets transport for HTTP client.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.transport = transport
	}
}
