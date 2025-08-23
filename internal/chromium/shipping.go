package chromium

import (
	"context"
	"fmt"
	"strings"

	addressvalidation "cloud.google.com/go/maps/addressvalidation/apiv1"
	"cloud.google.com/go/maps/addressvalidation/apiv1/addressvalidationpb"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/nyaruka/phonenumbers"
	"github.com/pistolricks/cosmetics-api/internal/validator"
	"google.golang.org/genproto/googleapis/type/postaladdress"
)

type lines []string

func ValidatePhone(v *validator.Validator, phone string) {
	v.Check(phone != "", "phone", "must be provided")
	v.Check(validator.Matches(phone, validator.PhoneRX), "phone", "must be a valid email address")
}

func (chrome ChromeClient) ProcessShipping(background func(fn func()), addressClient *addressvalidation.Client, email string, cookies []*proto.NetworkCookie, order goshopify.Order) bool {

	background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				chrome.InsertShippingInfo(addressClient, email, checkoutUrl, order)

			default:
				fmt.Println("not right cookie")
			}

		}
	})

	return true
}

func (chrome ChromeClient) InsertShippingInfo(addressClient *addressvalidation.Client, email string, checkoutUrl string, order goshopify.Order) {

	chrome.Client.Page.MustNavigate(checkoutUrl)

	shippingAddress := order.ShippingAddress

	firstName := strings.TrimSpace(shippingAddress.FirstName)
	lastName := strings.TrimSpace(shippingAddress.LastName)

	name := fmt.Sprintf("%s %s", firstName, lastName)

	address1 := strings.TrimSpace(shippingAddress.Address1)
	address2 := strings.TrimSpace(shippingAddress.Address2)
	company := strings.TrimSpace(shippingAddress.Company)
	city := strings.TrimSpace(shippingAddress.City)
	province := strings.TrimSpace(shippingAddress.Province)
	provinceCode := strings.TrimSpace(shippingAddress.ProvinceCode)
	shortZip := strings.TrimSpace(shippingAddress.Zip[:5])
	zip := strings.TrimSpace(shippingAddress.Zip)

	//phone := strings.Replace(strings.TrimSpace(shippingAddress.Phone), "+1", "", 1)

	num, err := phonenumbers.Parse(shippingAddress.Phone, "US")

	national := phonenumbers.Format(num, phonenumbers.E164)

	phone := strings.Replace(national, "+1", "", 1)

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

	postalAddress := postaladdress.PostalAddress{
		RegionCode:         shippingAddress.CountryCode,
		PostalCode:         shippingAddress.Zip,
		AdministrativeArea: shippingAddress.Province,
		Locality:           shippingAddress.City,
		Sublocality:        shippingAddress.City,
		AddressLines:       lines{shippingAddress.Address1, shippingAddress.Address2},
		Recipients:         lines{name, shippingAddress.Name},
		Organization:       shippingAddress.Company,
	}

	req := &addressvalidationpb.ValidateAddressRequest{Address: &postalAddress}
	resp, err := addressClient.ValidateAddress(context.Background(), req)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	location := resp.Result

	switch os := location.GetVerdict().AddressComplete; os {
	case true:

		address10 := location.Address.FormattedAddress

		address100 := location.GetUspsData().StandardizedAddress.FirstAddressLine

		var address20 string

		if company != "" && location.GetUspsData().StandardizedAddress.SecondAddressLine != "" {
			address20 = fmt.Sprintf("%s / %s", location.GetUspsData().StandardizedAddress.SecondAddressLine, company)
		} else if company == "" && location.GetUspsData().StandardizedAddress.SecondAddressLine == "" {
			address20 = ""
		} else if company != "" && location.GetUspsData().StandardizedAddress.SecondAddressLine == "" {
			address20 = fmt.Sprintf("%s", company)
		} else if company == "" && location.GetUspsData().StandardizedAddress.SecondAddressLine != "" {
			address20 = fmt.Sprintf("%s", location.GetUspsData().StandardizedAddress.SecondAddressLine)
		} else {
			address20 = fmt.Sprintf("%s", location.GetUspsData().StandardizedAddress.SecondAddressLine)
		}

		city0 := location.GetUspsData().StandardizedAddress.City
		// state0 := location.GetUspsData().StandardizedAddress.State
		zip0 := location.GetUspsData().StandardizedAddress.ZipCode
		chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(address10)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()
		chrome.Client.Page.MustWaitStable().MustElement("#address20").MustSelectAllText().MustInput(address20)
		chrome.Client.Page.MustWaitStable().MustElement("#city0").MustSelectAllText().MustInput(city0)
		// chrome.Client.Page.MustElement("#state0").MustSelect(state0)
		chrome.Client.Page.MustWaitStable().MustElement("#postalCode0").MustSelectAllText().MustInput(zip0)
		chrome.Client.Page.MustWaitStable().MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)

		chrome.Client.Page.MustWaitStable().MustElement("#address10").MustSelectAllText().MustInput(address100)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Escape).MustDo()

		chrome.Client.Page.MustWaitStable().MustElement("#email0").MustSelectAllText().MustInput(email)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()
	case false:
		address10 := address
		address20 := company
		city0 := location.GetUspsData().StandardizedAddress.City
		// state0 := location.GetUspsData().StandardizedAddress.State
		zip0 := location.GetUspsData().StandardizedAddress.ZipCode
		chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(address10)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()
		chrome.Client.Page.MustElement("#address20").MustSelectAllText().MustInput(address20)
		chrome.Client.Page.MustElement("#city0").MustSelectAllText().MustInput(city0)
		// chrome.Client.Page.MustElement("#state0").MustSelect(state0)
		chrome.Client.Page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip0)
		chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)
		chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

	default:
		address10 := location.Address.FormattedAddress

		address100 := location.GetUspsData().StandardizedAddress.FirstAddressLine

		var address20 string

		if company != "" && location.GetUspsData().StandardizedAddress.SecondAddressLine != "" {
			address20 = fmt.Sprintf("%s / %s", location.GetUspsData().StandardizedAddress.SecondAddressLine, company)
		} else if company == "" && location.GetUspsData().StandardizedAddress.SecondAddressLine == "" {
			address20 = ""
		} else if company != "" && location.GetUspsData().StandardizedAddress.SecondAddressLine == "" {
			address20 = fmt.Sprintf("%s", company)
		} else if company == "" && location.GetUspsData().StandardizedAddress.SecondAddressLine != "" {
			address20 = fmt.Sprintf("%s", location.GetUspsData().StandardizedAddress.SecondAddressLine)
		} else {
			address20 = fmt.Sprintf("%s", location.GetUspsData().StandardizedAddress.SecondAddressLine)
		}

		city0 := location.GetUspsData().StandardizedAddress.City
		// state0 := location.GetUspsData().StandardizedAddress.State
		zip0 := location.GetUspsData().StandardizedAddress.ZipCode
		chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(address10)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()
		chrome.Client.Page.MustWaitStable().MustElement("#address20").MustSelectAllText().MustInput(address20)
		chrome.Client.Page.MustWaitStable().MustElement("#city0").MustSelectAllText().MustInput(city0)
		// chrome.Client.Page.MustElement("#state0").MustSelect(state0)
		chrome.Client.Page.MustWaitStable().MustElement("#postalCode0").MustSelectAllText().MustInput(zip0)
		chrome.Client.Page.MustWaitStable().MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)

		chrome.Client.Page.MustWaitStable().MustElement("#address10").MustSelectAllText().MustInput(address100)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Escape).MustDo()

		chrome.Client.Page.MustWaitStable().MustElement("#email0").MustSelectAllText().MustInput(email)
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
		chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()
	}

	return
	/* Need to add Province/State */
	// chrome.Client.Page.MustElement("#state0").MustSelectAllText().MustInput(province)
}

func (chrome ChromeClient) ChromeShipping(address10 string, address20 string, city string, state string, zip string, phone string, email string) *rod.Element {

	chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(address10)

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()

	chrome.Client.Page.MustElement("#address20").MustSelectAllText().MustInput(address20)

	chrome.Client.Page.MustElement("#city0").MustSelectAllText().MustInput(city)
	chrome.Client.Page.MustElement("#state0").MustSelect(state)
	chrome.Client.Page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip)

	chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)

	chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

	return chrome.Client.Page.MustElement("#submitShipping")
}
