package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	v2 "github.com/pistolricks/cosmetics-api/internal/v2"
	log "github.com/sirupsen/logrus"
	"github.com/vinhluan/go-graphql-client"
	"github.com/vinhluan/go-shopify-graphql"
)

const (
	shopifyBaseDomain        = "myshopify.com"
	defaultAPIProtocol       = "https"
	defaultAPIBasePath       = "admin/api"
	defaultAPIEndpoint       = "graphql.json"
	defaultShopifyAPIVersion = "2025-07"
	defaultHttpTimeout       = time.Second * 10
)

type ClientApi struct {
	gql           graphql.GraphQL
	accessToken   string
	apiKey        string
	apiBasePath   string
	retries       int
	timeout       time.Duration
	transport     http.RoundTripper
	Product       *shopify.ProductService
	Variant       *shopify.VariantService
	Inventory     *shopify.InventoryService
	Collection    *shopify.CollectionService
	Order         *shopify.OrderService
	Fulfillment   *shopify.FulfillmentService
	Location      *shopify.LocationService
	Metafield     *shopify.MetafieldService
	BulkOperation *shopify.BulkOperationService
	Webhook       *shopify.WebhookService
}

func NewClient(shopName string, opts ...Option) *ClientApi {
	c := &ClientApi{
		apiBasePath: defaultAPIBasePath,
		timeout:     defaultHttpTimeout,
		transport:   http.DefaultTransport,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.gql == nil {
		apiEndpoint := buildAPIEndpoint(shopName, c.apiBasePath)
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

	c.Product = &v2.ProductV2{DB: db, client: c}
	c.Variant = &v2.VariantV2{DB: db, client client c, errors: k}
	c.Inventory = &InventoryV2{client: c}
	c.Collection = &CollectionV2{client: c}
	c.Order = &OrderV2{client: c}
	c.Fulfillment = &FulfillmentV2{client: c}
	c.Location = &LocationV2{client: c}
	c.Metafield = &MetafieldV2{client: c}
	c.BulkOperation = &BulkOperationV2{client: c}
	c.Webhook = &WebhookV2{client: c}

	return c
}

func (c ClientApi) NewDefaultClient() *ClientApi {
	apiKey := os.Getenv("STORE_API_KEY")
	accessToken := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	return NewClient(storeName, WithPrivateAppAuth(apiKey, accessToken), WithVersion(defaultShopifyAPIVersion))
}

func NewClientWithToken(accessToken string, storeName string) *ClientApi {
	if accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API access token and/or store name not set")
	}

	return NewClient(storeName, WithToken(accessToken), WithVersion(defaultShopifyAPIVersion))
}

func (c ClientApi) GraphQLClient() (graphql.GraphQL, graphql.GraphQL) {
	return c.gql, nil
}

func (c ClientApi) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.Mutate(ctx, m, variables)
		if err != nil {
			if r != nil {
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if IsConnectionError(err) {
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
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || IsConnectionError(err) {
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
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || IsConnectionError(err) {
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
