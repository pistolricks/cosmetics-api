package shopify

import (
	"errors"

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
