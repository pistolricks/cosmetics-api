package shopify

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

type FulfillClient struct {
	Config *ShopConfig
}

func (client FulfillClient) UpdateFulfillment(fulfillmentID uint64, trackingNumber string, trackingLink string) (*goshopify.Fulfillment, error) {

	updateFulfillment := goshopify.Fulfillment{
		Id: fulfillmentID,
		TrackingInfo: goshopify.FulfillmentTrackingInfo{
			Company: "Landmark Global",
			Number:  trackingNumber,
			Url:     trackingLink,
		},
		NotifyCustomer: true,
	}

	res, err := client.Config.Client.Order.UpdateFulfillment(context.Background(), fulfillmentID, updateFulfillment)

	fmt.Println("RES and ERR")
	fmt.Println(res)
	fmt.Println(err)
	fmt.Println(updateFulfillment)

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return res, nil
}
