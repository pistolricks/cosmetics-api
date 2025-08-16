package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pistolricks/cosmetics-api/internal/v2"
	log "github.com/sirupsen/logrus"
	"github.com/vinhluan/go-graphql-client"
)

const (
	shopifyBaseDomain = "myshopify.com"

	defaultAPIProtocol       = "https"
	defaultAPIBasePath       = "admin/api"
	defaultAPIEndpoint       = "graphql.json"
	defaultShopifyAPIVersion = "2023-04"
	defaultHttpTimeout       = time.Second * 10
)

// Product       v2.ProductService
// Variant       v2.VariantService
// Inventory     v2.InventoryService
// Collection    v2.CollectionService
// Order         v2.OrderService
// Fulfillment   v2.FulfillmentService
// Location      v2.LocationService
// Metafield     v2.MetafieldService
// BulkOperation v2.BulkOperationService
// Webhook       v2.WebhookService

type ClientApi struct {
	gql         graphql.GraphQL
	accessToken string
	apiKey      string
	apiBasePath string
	retries     int
	timeout     time.Duration
	transport   http.RoundTripper
}

func (c ClientApi) NewClient(shopName string, opts ...v2.Option) *ClientApi {
	c := ClientApi{
		apiBasePath: defaultAPIBasePath,
		timeout:     defaultHttpTimeout,
		transport:   http.DefaultTransport,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.gql == nil {
		apiEndpoint := v2.buildAPIEndpoint(shopName, c.apiBasePath)
		httpClient := &http.Client{
			Timeout: c.timeout,
			Transport: &transport{
				accessToken:  c.accessToken,
				apiKey:       c.apiKey,
				apiBasePath:  c.apiBasePath,
				roundTripper: c.transport,
			},
		}
		c.gql = graphql.NewClient(apiEndpoint, httpClient)
	}

	return c
}

func (c ClientApi) NewDefaultClient() *ClientApi {
	apiKey := os.Getenv("STORE_API_KEY")
	accessToken := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	return c.NewClient(storeName, v2.WithPrivateAppAuth(apiKey, accessToken), v2.WithVersion(defaultShopifyAPIVersion))
}

func (c ClientApi) NewClientWithToken(accessToken string, storeName string) *ClientApi {
	if accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API access token and/or store name not set")
	}

	return c.NewClient(storeName, v2.WithToken(accessToken), v2.WithVersion(defaultShopifyAPIVersion))
}

func (c ClientApi) GraphQLClient() graphql.GraphQL {
	return c.gql
}

func (c ClientApi) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.Mutate(ctx, m, variables)
		if err != nil {
			if r != nil {
				wait := v2.CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if v2.IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}

func (c ClientApi) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.Query(ctx, q, variables)
		if err != nil {
			if r != nil {
				wait := v2.CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || v2.IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}

func (c ClientApi) QueryString(ctx context.Context, q string, variables map[string]interface{}, out interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.QueryString(ctx, q, variables, out)
		if err != nil {
			if r != nil {
				wait := v2.CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || v2.IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}
