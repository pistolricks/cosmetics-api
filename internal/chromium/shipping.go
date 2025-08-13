package chromium

import (
	"fmt"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func (chrome ChromeClient) ProcessShipping(background func(fn func()), email string, cookies []*proto.NetworkCookie, order goshopify.Order) {

	background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				chrome.InsertShippingInfo(email, checkoutUrl, order)

			default:
				fmt.Println("not right cookie")
			}

		}
	})

}

func (chrome ChromeClient) InsertShippingInfo(email string, checkoutUrl string, order goshopify.Order) {

	chrome.Client.Page.MustNavigate(checkoutUrl)

	shippingAddress := order.ShippingAddress

	firstName := strings.TrimSpace(shippingAddress.FirstName)
	lastName := strings.TrimSpace(shippingAddress.LastName)

	address1 := strings.TrimSpace(shippingAddress.Address1)
	address2 := strings.TrimSpace(shippingAddress.Address2)
	company := strings.TrimSpace(shippingAddress.Company)
	city := strings.TrimSpace(shippingAddress.City)
	province := strings.TrimSpace(shippingAddress.Province)
	provinceCode := strings.TrimSpace(shippingAddress.ProvinceCode)
	shortZip := strings.TrimSpace(shippingAddress.Zip[:5])
	zip := strings.TrimSpace(shippingAddress.Zip)

	phone := strings.Replace(strings.TrimSpace(shippingAddress.Phone), "+1", "", 1)
	// email := strings.TrimSpace(order.Email)

	chrome.Client.Page.MustElement("#firstName0").MustSelectAllText().MustInput(firstName)
	chrome.Client.Page.MustElement("#lastName0").MustSelectAllText().MustInput(lastName)

	removedAddress2 := strings.Replace(address1, address2, "", 1)
	removedCity := strings.Replace(removedAddress2, city, "", 1)
	removedProvince := strings.Replace(removedCity, province, "", 1)
	removedProvinceCode := strings.Replace(removedProvince, provinceCode, "", 1)
	removedZip := strings.Replace(removedProvinceCode, zip, "", 1)
	lineAddress := strings.Replace(removedZip, shortZip, "", 1)

	formattedAddress := strings.TrimSpace(lineAddress)

	address := fmt.Sprintf("%s %s, %s", formattedAddress, address2, shortZip)

	chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(address)

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()

	chrome.Client.Page.MustElement("#address20").MustSelectAllText().MustInput(company)

	// chrome.Client.Page.MustElement("#city0").MustSelectAllText().MustInput(city)
	// chrome.Client.Page.MustElement("#state0").MustSelect(provinceCode)
	// chrome.Client.Page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip)

	chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)

	chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

	/* Need to add Province/State */
	// chrome.Client.Page.MustElement("#state0").MustSelectAllText().MustInput(province)
}
