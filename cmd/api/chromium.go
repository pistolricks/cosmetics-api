package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"golang.org/x/net/context"
)

func (app *application) RimanLogin(loginUrl string, rimanStoreName string, username string, password string) (*rod.Page, *rod.Browser, []*proto.NetworkCookie) {

	// --allow-third-party-cookies
	path, _ := launcher.LookPath()

	officeUrl := fmt.Sprintf("https://myoffice.riman.com")
	homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	u := launcher.
		NewUserMode().
		UserDataDir("path").
		Headless(false).
		NoSandbox(true).
		Bin(path).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect().NoDefaultDevice()

	page := browser.MustPage().MustWindowNormal()

	wait := page.MustWaitNavigation()
	page.MustNavigate(officeUrl)
	wait()

	url := page.MustInfo().URL

	switch {
	case strings.Contains(url, "https://myoffice-1.riman.com/login"):
		fmt.Println("office login")

		page.MustElement("#username").MustInput(username)
		page.MustElement("#password").MustInput(password)
		page.MustElement("#loginBtn").MustClick()

	case strings.Contains(url, officeUrl):
		fmt.Println("office logged in")
		page.MustNavigate(homeUrl)
	}

	cookies := browser.MustGetCookies()
	return page, browser, cookies
}

func (app *application) HomePage(rimanStoreName string, page *rod.Page, browser *rod.Browser) (*rod.Page, *rod.Browser, []*proto.NetworkCookie) {
	// networkCookie := networkCookies(cookies)

	homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	// page.MustSetCookies(networkCookie...)

	wait := page.MustWaitNavigation()
	page.MustNavigate(homeUrl)
	wait()

	app.page = page
	app.browser = browser

	cookies := browser.MustGetCookies()

	return page, browser, cookies
}

func (app *application) ProcessOrders(rimanStoreName string, page *rod.Page, browser *rod.Browser, cookies []*proto.NetworkCookie, orders []goshopify.Order) {
	orderCount := len(orders)

	switch orderCount := orderCount; {
	case orderCount == 1:
		app.SubmitOrder(rimanStoreName, page, browser, cookies, orders[0])
	case orderCount > 1:
		for _, order := range orders {
			app.SubmitOrder(rimanStoreName, page, browser, cookies, order)
		}
	}
}

func (app *application) SubmitOrder(rimanStoreName string, page *rod.Page, browser *rod.Browser, cookies []*proto.NetworkCookie, order goshopify.Order) {

	count := len(order.LineItems)

	for i, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/%s/products/%s", rimanStoreName, product.SKU)

		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl) // := browser.MustPage(productUrl)
		wait()

		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))
		page.MustElement("button.add-to-bag-btn").MustClick()
		page.MustWaitStable()

		fmt.Println("QTY, and Index")
		println(i + 1)
		println(count)
		/*
			switch {
			case i < count-1:
				page.MustElement("div.cart-btn").MustClick()
			case i == count-1:

				cookies, err := browser.GetCookies()
				if err != nil {
					return
				}

				app.processShipping(browser, page, cookies, order)

			}
		*/

		app.processShipping(browser, page, cookies, order)
	}
}

func networkCookies(cookies []*proto.NetworkCookie) []*proto.NetworkCookieParam {

	var networkCookie []*proto.NetworkCookieParam

	for _, cookie := range cookies {
		networkCookie = append(networkCookie, &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
			SameSite: cookie.SameSite,
			Expires:  cookie.Expires,
		})
	}

	return networkCookie
}

func (app *application) processShipping(browser *rod.Browser, page *rod.Page, cookies []*proto.NetworkCookie, order goshopify.Order) {

	app.background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				app.insertShippingInfo(browser, page, checkoutUrl, order)

			default:
				fmt.Println("not right cookie")
			}

		}
	})

}

type StateObject = struct {
	Code  string
	Name  string
	name2 any
}

/* TODO: REMOVE HARD CODED EMAIL */

func IsValidPhoneNumber(phone_number string) bool {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phone_number = strings.ReplaceAll(phone_number, " ", "")

	return re.Find([]byte(phone_number)) != nil
}

func (app *application) insertShippingInfo(browser *rod.Browser, page *rod.Page, checkoutUrl string, order goshopify.Order) {

	page.MustNavigate(checkoutUrl)

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

	page.MustElement("#firstName0").MustSelectAllText().MustInput(firstName)
	page.MustElement("#lastName0").MustSelectAllText().MustInput(lastName)

	removedAddress2 := strings.Replace(address1, address2, "", 1)
	removedCity := strings.Replace(removedAddress2, city, "", 1)
	removedProvince := strings.Replace(removedCity, province, "", 1)
	removedProvinceCode := strings.Replace(removedProvince, provinceCode, "", 1)
	removedZip := strings.Replace(removedProvinceCode, zip, "", 1)
	lineAddress := strings.Replace(removedZip, shortZip, "", 1)

	formattedAddress := strings.TrimSpace(lineAddress)

	address := fmt.Sprintf("%s %s, %s", formattedAddress, address2, shortZip)

	page.MustElement("#address10").MustSelectAllText().MustInput(address)
	page.MustElement("#address20").MustSelectAllText().MustInput(company)

	page.MustElement("#city0").MustSelectAllText().MustInput(city)
	// page.MustElement("#state0").MustSelect(provinceCode)
	page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip)

	page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)
	email := app.client.Email
	page.MustElement("#email0").MustSelectAllText().MustInput(email)

	/* Need to add Province/State */
	// page.MustElement("#state0").MustSelectAllText().MustInput(province)
}

func (app *application) insertBillingInfo(page *rod.Page) bool {

	// nameOnCard := os.Getenv("NAME_ON_CARD")
	// cardNumber := os.Getenv("CARD_NUMBER")
	// cardMonth := os.Getenv("CARD_MONTH")
	// cardYear := os.Getenv("CARD_YEAR")
	// cardCvv := os.Getenv("CARD_CVV")

	billingFirstName := os.Getenv("BILLING_FIRST_NAME")
	billingLastName := os.Getenv("BILLING_LAST_NAME")
	billingAddress1 := os.Getenv("BILLING_ADDRESS_1")
	// billingAddress2 := os.Getenv("BILLING_ADDRESS_2")
	billingCity := os.Getenv("BILLING_CITY")
	// billingState := os.Getenv("BILLING_STATE")
	billingZip := os.Getenv("BILLING_ZIP")
	billingPhone := os.Getenv("BILLING_PHONE")

	el := page.MustElement("#mat-checkbox-7-input")

	// check if it is checked
	if el.MustProperty("checked").Bool() {
		el.MustClick()
	}

	page.MustElement("#firstName0").MustSelectAllText().MustInput(billingFirstName)
	page.MustElement("#lastName0").MustSelectAllText().MustInput(billingLastName)

	// Name on Card
	page.MustElement("#address10").MustSelectAllText().MustInput(billingAddress1)
	// page.MustElement("#address20").MustSelectAllText().MustInput(billingAddress2)

	page.MustElement("#city0").MustSelectAllText().MustInput(billingCity)
	// page.MustElement("#state0").MustSelect(provinceCode)
	_ = page.MustElement("#state0").Select([]string{`[value="California"]`}, true, rod.SelectorTypeCSSSector)

	page.MustElement("#postalCode0").MustSelectAllText().MustInput(billingZip)

	page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(billingPhone)
	email := app.client.Email
	page.MustElement("#email0").MustSelectAllText().MustInput(email)

	/* Need to add Province/State */
	// page.MustElement("#state0").MustSelectAllText().MustInput(province)

	return true
}

func (app *application) emulateGPS(page *rod.Page) bool {
	latStr := os.Getenv("LAT")
	lngStr := os.Getenv("LNG")

	// Convert string coordinates to float64
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
	}

	accuracy := 100.0 // Define accuracy in meters

	override := &proto.EmulationSetGeolocationOverride{
		Latitude:  &lat,
		Longitude: &lng,
		Accuracy:  &accuracy, // Accuracy is optional but recommended.
	}

	_, err = page.Call(context.Background(), "", "Emulation.setGeolocationOverride", override)
	if err != nil {
		fmt.Println("Failed to set geolocation:", err)
	}

	return true
}
