package shopify

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"

	"github.com/vinhluan/go-shopify-graphql/model"
)

type OrderClient struct {
	Config *ShopConfig
}

/*
https://cart-api.riman.com/api/v2/order
{
    "mainId": 47387,
    "mainOrderType": 4,
    "countryCode": "US",
    "salesCampaignFK": null,
    "cartKey": "7fe33b2a-b99c-4244-9707-b397bea92eff"
}
*/

func (orderClient OrderClient) UpdateOrderNote(orderId uint64, note string) (*goshopify.Order, error) {

	o := goshopify.Order{
		Id:   orderId, // orderId,
		Note: note,
	}

	order, err := orderClient.Config.Client.Order.Update(context.Background(), o)
	if err != nil {
		return order, err
	}

	expected := goshopify.Order{Id: orderId}
	if o.Id != expected.Id {
		return order, err
	}

	return order, err
}

func ListAllOrders() ([]*model.Order, error) {

	return nil, nil
}
