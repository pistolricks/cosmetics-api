package chromium

import (
	"os"

	"github.com/go-rod/rod/lib/input"
)

func (chrome ChromeClient) InsertBillingInfo(email string) bool {

	//chrome.Client.Page.Timeout(2 * time.Second).MustNavigate("https://mall.riman.com/checkout/billing?orderNum=0x0200000025227197FB7D1FAD7A856026F088D63E312861E9D6C8DF1D29F8C4327B40D24E7E35573BC1B1190A08488CD351124AF1")
	chrome.Client.Page.MustActivate()

	chrome.Client.Page.MustElement("#firstName0").MustSelectAllText().MustInput(os.Getenv("BILLING_FIRST_NAME"))
	chrome.Client.Page.MustElement("#lastName0").MustSelectAllText().MustInput(os.Getenv("BILLING_LAST_NAME"))

	chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(os.Getenv("BILLING_ADDRESS_1"))

	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.ArrowDown).MustDo()
	chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Enter).MustDo()

	chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(os.Getenv("BILLING_PHONE"))
	chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)

	return true

}

// chrome.Client.Page.MustElement("#mat-input-10").MustKeyActions().Type(input.Tab).MustDo()
// chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()
// chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Space).MustDo()
// chrome.Client.Page.MustWaitStable().KeyActions().Type(input.Tab).MustDo()

/*
	nameOnCard := os.Getenv("NAME_ON_CARD")
	cardNumber := os.Getenv("CARD_NUMBER")
	cardMonth := os.Getenv("CARD_MONTH")
	cardYear := os.Getenv("CARD_YEAR")
	cardCvv := os.Getenv("CARD_CVV")
*/

// billingFirstName := os
// billingLastName := os
// billingAddress1 := os.Getenv("BILLING_ADDRESS_1")
// billingAddress2 := os.Getenv("BILLING_ADDRESS_2")
// billingCity := os.Getenv("BILLING_CITY")
// billingState := os.Getenv("BILLING_STATE")
// billingZip := os.Getenv("BILLING_ZIP")
// billingPhone := os.Getenv("BILLING_PHONE")

// chrome.Client.Page.MustElement("#firstName0").MustSelectAllText().MustInput(billingFirstName)
// chrome.Client.Page.MustElement("#lastName0").MustSelectAllText().MustInput(billingLastName)
// chrome.Client.Page.MustElement("#address10").MustSelectAllText().MustInput(billingAddress1)
// chrome.Client.Page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(billingPhone)
//
// chrome.Client.Page.MustElement("#email0").MustSelectAllText().MustInput(email)

type Billing struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	ShipFName interface{} `json:"shipFName"`
	ShipLName interface{} `json:"shipLName"`
	Address1  string      `json:"address1"`
	Address2  string      `json:"address2"`
	Address3  string      `json:"address3"`
	City      string      `json:"city"`
	CityName  interface{} `json:"cityName"`
	Zip       string      `json:"zip"`
	State     struct {
		Code  string      `json:"code"`
		Name  string      `json:"name"`
		Name2 interface{} `json:"name2"`
	} `json:"state"`
	Phone       string `json:"phone"`
	SecondPhone string `json:"secondPhone"`
	Email       string `json:"email"`
	Country     struct {
		Code2  string      `json:"code2"`
		States interface{} `json:"states"`
	} `json:"country"`
	Ssn                  interface{} `json:"ssn"`
	Area                 string      `json:"area"`
	AreaName             interface{} `json:"areaName"`
	SiteUrl              interface{} `json:"siteUrl"`
	IsUseShippingAddress bool        `json:"isUseShippingAddress"`
}
