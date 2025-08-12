package riman

import (
	"fmt"
)

type Cart struct {
	CartKey                   string      `json:"cartKey"`
	CartType                  string      `json:"cartType"`
	CountryCode               string      `json:"countryCode"`
	MainFK                    int         `json:"mainFK"`
	MainReferrerFK            int         `json:"mainReferrerFK"`
	Culture                   string      `json:"culture"`
	LanguageFK                int         `json:"languageFK"`
	MainOrderTypeFK           int         `json:"mainOrderTypeFK"`
	PriceListFK               int         `json:"priceListFK"`
	CampaignCode              string      `json:"campaignCode"`
	SalesCampaignFK           interface{} `json:"salesCampaignFK"`
	Ip                        string      `json:"ip"`
	DateEntered               string      `json:"dateEntered"`
	GaCode                    string      `json:"gaCode"`
	FacebookCode              string      `json:"facebookCode"`
	LuckyOrange               string      `json:"luckyOrange"`
	ReferrerSiteUrl           string      `json:"referrerSiteUrl"`
	ReferrerIsCorporate       bool        `json:"referrerIsCorporate"`
	CustomerReferralId        string      `json:"customerReferralId"`
	MainCreditCardsFK         interface{} `json:"mainCreditCardsFK"`
	MainOrdersFK              interface{} `json:"mainOrdersFK"`
	ShippingTypeFK            int         `json:"shippingTypeFK"`
	CartStatus                interface{} `json:"cartStatus"`
	FirstName                 string      `json:"firstName"`
	LastName                  string      `json:"lastName"`
	Phone                     string      `json:"phone"`
	Email                     string      `json:"email"`
	DateModified              string      `json:"dateModified"`
	SubTotal                  int         `json:"subTotal"`
	FormattedSubTotal         string      `json:"formattedSubTotal"`
	Tax                       int         `json:"tax"`
	FormattedTax              string      `json:"formattedTax"`
	Shipping                  int         `json:"shipping"`
	FormattedShipping         string      `json:"formattedShipping"`
	Discount                  int         `json:"discount"`
	FormattedDiscount         string      `json:"formattedDiscount"`
	Total                     int         `json:"total"`
	FormattedTotal            string      `json:"formattedTotal"`
	PointsTotal               float64     `json:"pointsTotal"`
	ShipSignatureRequired     bool        `json:"shipSignatureRequired"`
	ShipSignatureFee          int         `json:"shipSignatureFee"`
	CurrencyFK                int         `json:"currencyFK"`
	CurrencyCode              string      `json:"currencyCode"`
	MainDiscountCode          string      `json:"mainDiscountCode"`
	ActiveSmartDelivery       bool        `json:"activeSmartDelivery"`
	AllowImportCart           bool        `json:"allowImportCart"`
	OfferPreferredCust        bool        `json:"offerPreferredCust"`
	IsAffiliateOn             bool        `json:"isAffiliateOn"`
	IsVolumeBasedRSB          bool        `json:"isVolumeBasedRSB"`
	OfferLoyaltyProgram       bool        `json:"offerLoyaltyProgram"`
	OfferSmartDelivery        bool        `json:"offerSmartDelivery"`
	IsRetailSignup            bool        `json:"isRetailSignup"`
	HasRetailStarterKit       bool        `json:"hasRetailStarterKit"`
	AllowRetailSignup         bool        `json:"allowRetailSignup"`
	EventMemberID             interface{} `json:"eventMemberID"`
	ShowAbandonedOrderWarning bool        `json:"showAbandonedOrderWarning"`
	ShouldCreateAccount       bool        `json:"shouldCreateAccount"`
	IsCartCreatedFromBag      bool        `json:"isCartCreatedFromBag"`
	IsCartCreatedFromSignup   bool        `json:"isCartCreatedFromSignup"`
	CurrencySymbol            string      `json:"currencySymbol"`
	ShippingAddress           struct {
		MainCartAddressPK int         `json:"mainCartAddressPK"`
		FirstName         interface{} `json:"firstName"`
		LastName          interface{} `json:"lastName"`
		Company           interface{} `json:"company"`
		Address1          interface{} `json:"address1"`
		Address2          interface{} `json:"address2"`
		Address3          interface{} `json:"address3"`
		City              interface{} `json:"city"`
		State             interface{} `json:"state"`
		PostalCode        interface{} `json:"postalCode"`
		Country           interface{} `json:"country"`
		IsPOBox           interface{} `json:"isPOBox"`
		IsResidential     interface{} `json:"isResidential"`
	} `json:"shippingAddress"`
	MailingAddress struct {
		MainCartAddressPK int         `json:"mainCartAddressPK"`
		FirstName         interface{} `json:"firstName"`
		LastName          interface{} `json:"lastName"`
		Company           interface{} `json:"company"`
		Address1          interface{} `json:"address1"`
		Address2          interface{} `json:"address2"`
		Address3          interface{} `json:"address3"`
		City              interface{} `json:"city"`
		State             interface{} `json:"state"`
		PostalCode        interface{} `json:"postalCode"`
		Country           interface{} `json:"country"`
		IsPOBox           interface{} `json:"isPOBox"`
		IsResidential     interface{} `json:"isResidential"`
	} `json:"mailingAddress"`
	BillingAddress struct {
		MainCartAddressPK int         `json:"mainCartAddressPK"`
		FirstName         interface{} `json:"firstName"`
		LastName          interface{} `json:"lastName"`
		Company           interface{} `json:"company"`
		Address1          interface{} `json:"address1"`
		Address2          interface{} `json:"address2"`
		Address3          interface{} `json:"address3"`
		City              interface{} `json:"city"`
		State             interface{} `json:"state"`
		PostalCode        interface{} `json:"postalCode"`
		Country           interface{} `json:"country"`
		IsPOBox           interface{} `json:"isPOBox"`
		IsResidential     interface{} `json:"isResidential"`
	} `json:"billingAddress"`
	FormattedAutoshipSubtotal string `json:"formattedAutoshipSubtotal"`
	CartItems                 []struct {
		Id                          int           `json:"id"`
		Quantity                    int           `json:"quantity"`
		ProductFk                   int           `json:"productFk"`
		PackageFk                   interface{}   `json:"packageFk"`
		Name                        string        `json:"name"`
		ImageUrl                    string        `json:"imageUrl"`
		SetupForAs                  bool          `json:"setupForAs"`
		ConfigFk                    interface{}   `json:"configFk"`
		PriceListFk                 int           `json:"priceListFk"`
		IsPaCItem                   bool          `json:"isPaCItem"`
		IsSignup                    bool          `json:"isSignup"`
		Sku                         string        `json:"sku"`
		PriceType                   string        `json:"priceType"`
		CountryCode                 string        `json:"countryCode"`
		CurrencyPK                  int           `json:"currencyPK"`
		CurrencyCode                string        `json:"currencyCode"`
		CurrencySymbol              string        `json:"currencySymbol"`
		ExtraFee                    interface{}   `json:"extraFee"`
		Discount                    int           `json:"discount"`
		FormattedDiscount           string        `json:"formattedDiscount"`
		BasePrice                   int           `json:"basePrice"`
		FormattedBasePrice          string        `json:"formattedBasePrice"`
		UnitPrice                   int           `json:"unitPrice"`
		FormattedUnitPrice          string        `json:"formattedUnitPrice"`
		RetailTaxablePrice          int           `json:"retailTaxablePrice"`
		FormattedRetailTaxablePrice string        `json:"formattedRetailTaxablePrice"`
		Cv                          int           `json:"cv"`
		Sp                          int           `json:"sp"`
		MaxLimit                    int           `json:"maxLimit"`
		Points                      float64       `json:"points"`
		IsShippable                 bool          `json:"isShippable"`
		IsStarterKit                bool          `json:"isStarterKit"`
		BrandName                   string        `json:"brandName"`
		LineSubTotal                int           `json:"lineSubTotal"`
		FormattedLineSubTotal       string        `json:"formattedLineSubTotal"`
		LineTotal                   int           `json:"lineTotal"`
		FormattedLineTotal          string        `json:"formattedLineTotal"`
		OfferAffiliateProgram       bool          `json:"offerAffiliateProgram"`
		IsVolumeBasedRSB            bool          `json:"isVolumeBasedRSB"`
		OfferPreferredCust          bool          `json:"offerPreferredCust"`
		OfferLoyaltyProgram         bool          `json:"offerLoyaltyProgram"`
		ShowSDCheckbox              bool          `json:"showSDCheckbox"`
		JoinMaxLifetimeLimitCatCode string        `json:"joinMaxLifetimeLimitCatCode"`
		JoinMaxLifetimeLimit        int           `json:"joinMaxLifetimeLimit"`
		ChildItems                  []interface{} `json:"childItems"`
	} `json:"cartItems"`
}

type CartErrors struct {
	Error string `json:"error"`
}

func GetCart(cartKey string) (*Cart, error) {

	cartUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s", cartKey)

	fmt.Println(cartUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&Cart{}).
		SetError(&CartErrors{}).
		Get(cartUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*Cart))

	return res.Result().(*Cart), err
}
