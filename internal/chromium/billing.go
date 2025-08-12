package chromium

import (
	"os"

	"github.com/go-rod/rod"
)

func (chrome ChromeClient) InsertBillingInfo(email string) bool {

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

	el := chrome.Client.Page.MustElement("#mat-checkbox-7-input")

	// check if it is checked
	if el.MustProperty("checked").Bool() {
		el.MustClick()
	}

	chrome.Client.Page.MustElement("#firstName0").MustSelectAllText().MustInput(billingFirstName)
	chrome.Client.Page.MustElement("#lastName0").MustSelectAllText().MustInput(billingLastName)

	// Name on Card
	chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(billingAddress1)
	// chrome.Client.Page.MustElement("#address20").MustSelectAllText().MustInput(billingAddress2)

	chrome.Client.Page.MustElement("#city0").MustSelectAllText().MustInput(billingCity)
	// chrome.Client.Page.MustElement("#state0").MustSelect(provinceCode)
	_ = chrome.Client.Page.MustElement("#state0").Select([]string{`[value="California"]`}, true, rod.SelectorTypeCSSSector)

	chrome.Client.Page.MustElement("#postalCode0").MustSelectAllText().MustInput(billingZip)

	chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(billingPhone)

	chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)

	/* Need to add Province/State */
	// page.MustElement("#state0").MustSelectAllText().MustInput(province)

	return true
}
