package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

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

type ListOptions struct {
	Query   string
	First   int
	Last    int
	After   string
	Before  string
	Reverse bool
}
type Client struct {
	gql            graphql.GraphQL
	accessToken    string
	apiKey         string
	apiBasePath    string
	retries        int
	timeout        time.Duration
	transport      http.RoundTripper
	Products       ProductService
	Variants       VariantService
	Inventory      InventoryService
	Collections    CollectionService
	Orders         OrderService
	Fulfillments   FulfillmentService
	Locations      LocationService
	Metafields     MetafieldService
	BulkOperations BulkOperationService
	Webhooks       WebhookService
}

func NewClient(shopName string, opts ...Option) *Client {
	c := &Client{
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

	c.Products = &ProductServiceOp{client: c}
	c.Variants = &VariantServiceOp{client: c}
	c.Inventorys = &InventoryServiceOp{client: c}
	c.Collections = &CollectionServiceOp{client: c}
	c.Orders = &OrderServiceOp{client: c}
	c.Fulfillments = &FulfillmentServiceOp{client: c}
	c.Locations = &LocationServiceOp{client: c}
	c.Metafields = &MetafieldServiceOp{client: c}
	c.BulkOperations = &BulkOperationServiceOp{client: c}
	c.Webhooks = &WebhookServiceOp{client: c}

	return c
}

func NewDefaultClient() *Client {
	apiKey := os.Getenv("STORE_API_KEY")
	accessToken := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	return NewClient(storeName, WithPrivateAppAuth(apiKey, accessToken), WithVersion(defaultShopifyAPIVersion))
}

func NewClientWithToken(accessToken string, storeName string) *Client {
	if accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API access token and/or store name not set")
	}

	return NewClient(storeName, WithToken(accessToken), WithVersion(defaultShopifyAPIVersion))
}

func (c *Client) GraphQLClient() graphql.GraphQL {
	return c.gql
}

func (c *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
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

func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
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

func (c *Client) QueryString(ctx context.Context, q string, variables map[string]interface{}, out interface{}) error {
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
