package main

import (
	"net/http"
	"os"

	graphify "github.com/vinhluan/go-graphql-client"
)

func (app *application) NewClient(shopName string, opts ...Option) *Client {
	c := &Client{
		apiBasePath: app.config.graphql.defaultAPIBasePath,
		timeout:     app.config.graphql.defaultHTTPTimeout,
		transport:   http.DefaultTransport,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.gql == nil {
		apiEndpoint := app.buildAPIEndpoint(shopName, c.apiBasePath)
		httpClient := &http.Client{
			Timeout: c.timeout,
			Transport: Transport{
				accessToken:  c.accessToken,
				apiKey:       c.apiKey,
				apiBasePath:  c.apiBasePath,
				roundTripper: c.transport,
			},
		}
		c.gql = graphify.NewClient(apiEndpoint, httpClient)
	}

	//c.Products = &ProductServiceOp{client: c}
	//c.Variants = &VariantServiceOp{client: c}
	//c.Inventory = &InventoryerviceOp{client: c}
	//c.Collections = &CollectionServiceOp{client: c}
	//c.Orders = &OrderServiceOp{client: c}
	//c.Fulfillments = &FulfillmentServiceOp{client: c}
	//c.Locations = &LocationServiceOp{client: c}
	//c.Metafields = &MetafieldServiceOp{client: c}
	//c.BulkOperations = &BulkOperationServiceOp{client: c}
	//c.Webhooks = &WebhookServiceOp{client: c}

	return c
}

func (app *application) NewDefaultClient() *Client {
	apiKey := os.Getenv("STORE_API_KEY")
	accessToken := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || accessToken == "" || storeName == "" {
		app.logError("ERROR", "Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	return NewClient(storeName, WithPrivateAppAuth(apiKey, accessToken), WithVersion(defaultShopifyAPIVersion))
}

func (app *application) NewClientWithToken(accessToken string, storeName string) *Client {
	if accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API access token and/or store name not set")
	}

	return NewClient(storeName, WithToken(accessToken), WithVersion(defaultShopifyAPIVersion))
}
