package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

type ShopConfig struct {
	App       *goshopify.App
	Client    *goshopify.Client
	ShopName  string
	ShopToken string
}

type ShopClient struct {
	Orders OrderClient
}

func NewShopClient(config *ShopConfig) ShopClient {
	return ShopClient{
		Orders: OrderClient{Config: config},
	}
}
