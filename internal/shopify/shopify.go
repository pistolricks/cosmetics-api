package shopify

import (
	"errors"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

var (
	ErrShopifyApi = errors.New("shopify API error")
)

type ShopConfig struct {
	OrderApp  *goshopify.App
	Client    *goshopify.Client
	ShopName  string
	ShopToken string
}

type FulfillmentClient struct {
	Config *ShopConfig
}

type ShopClient struct {
	Orders       OrderClient
	Fulfillments FulfillmentClient
}

func NewShopClient(config ShopConfig) ShopClient {
	return ShopClient{
		Orders:       OrderClient{Config: &config},
		Fulfillments: FulfillmentClient{Config: &config},
	}
}

func ShopifyV1() ShopConfig {
	app := goshopify.App{
		ApiKey:      os.Getenv("SHOPIFY_KEY"),
		ApiSecret:   os.Getenv("SHOPIFY_SECRET"),
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders,write_orders,read_fulfillments,write_fulfillments",
	}

	client, err := goshopify.NewClient(app, os.Getenv("STORE_NAME"), os.Getenv("SHOPIFY_TOKEN"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	shopConfig := ShopConfig{OrderApp: &app, Client: client, ShopName: os.Getenv("STORE_NAME"), ShopToken: os.Getenv("SHOPIFY_TOKEN")}

	return shopConfig
}
