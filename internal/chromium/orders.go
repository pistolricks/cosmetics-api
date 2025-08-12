package chromium

import (
	"fmt"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod/lib/proto"
)

func (chrome ChromeClient) ProcessOrders(background func(fn func()), email, rimanStoreName string, cookies []*proto.NetworkCookie, orders []goshopify.Order) {
	orderCount := len(orders)

	switch orderCount := orderCount; {
	case orderCount == 1:
		chrome.SubmitOrder(background, email, rimanStoreName, cookies, orders[0])
	case orderCount > 1:
		for _, order := range orders {
			chrome.SubmitOrder(background, email, rimanStoreName, cookies, order)
		}
	}
}

func (chrome ChromeClient) SubmitOrder(background func(fn func()), email, rimanStoreName string, cookies []*proto.NetworkCookie, order goshopify.Order) {

	count := len(order.LineItems)

	for i, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/%s/products/%s", rimanStoreName, product.SKU)

		wait := chrome.Client.Page.MustWaitNavigation()
		chrome.Client.Page.MustNavigate(productUrl) // := browser.MustPage(productUrl)
		wait()

		chrome.Client.Page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))
		chrome.Client.Page.MustElement("button.add-to-bag-btn").MustClick()
		chrome.Client.Page.MustWaitStable()

		fmt.Println("QTY, and Index")
		println(i + 1)
		println(count)

		chrome.ProcessShipping(background, email, cookies, order)
	}
}
