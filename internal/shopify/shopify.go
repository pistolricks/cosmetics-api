package shopify

import (
	"errors"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

var (
	ErrShopifyApi = errors.New("shopify API error")
)

type ShopConfig struct {
	OrderApp   *goshopify.App
	FulfillApp *goshopify.App
	Client     *goshopify.Client
	ShopName   string
	ShopToken  string
}

type ShopClient struct {
	Orders       OrderClient
	Fulfillments FulfillClient
}

func NewShopClient(config *ShopConfig) ShopClient {

	return ShopClient{
		Orders:       OrderClient{Config: config},
		Fulfillments: FulfillClient{Config: config},
	}

}

func NewShopifyClient(storeName string, shopifyToken string) (*goshopify.Client, error) {
	sa := goshopify.App{
		ApiKey:      os.Getenv("SHOPIFY_KEY"),
		ApiSecret:   os.Getenv("SHOPIFY_SECRET"),
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders,write_orders,read_fulfillments,write_fulfillments,write_third_party_fulfillment_orders,read_shipping,write_shipping",
	}

	client, err := goshopify.NewClient(sa, storeName, shopifyToken)
	if err != nil {
		err.Error()
		os.Exit(1)
	}

	return client, nil
}

func ShopifyConfig(client *goshopify.Client) ShopConfig {
	sa := goshopify.App{
		ApiKey:      os.Getenv("SHOPIFY_KEY"),
		ApiSecret:   os.Getenv("SHOPIFY_SECRET"),
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders,write_orders,read_fulfillments,write_fulfillments,write_third_party_fulfillment_orders,read_shipping,write_shipping",
	}

	return ShopConfig{OrderApp: &sa, Client: client, ShopName: os.Getenv("STORE_NAME"), ShopToken: os.Getenv("SHOPIFY_TOKEN")}
}

func OrderConfig(client *goshopify.Client) ShopConfig {

	sa := goshopify.App{
		ApiKey:      os.Getenv("SHOPIFY_KEY"),
		ApiSecret:   os.Getenv("SHOPIFY_SECRET"),
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders,write_orders,read_fulfillments,write_fulfillments,write_third_party_fulfillment_orders,read_shipping,write_shipping",
	}
	return ShopConfig{OrderApp: sa, Client: client, ShopName: os.Getenv("STORE_NAME"), ShopToken: os.Getenv("SHOPIFY_TOKEN")}
}
