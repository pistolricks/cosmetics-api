package riman

type RimanProducts struct {
	IsProductOnAutoship          bool            `json:"isProductOnAutoship"`
	SeqNo                        int             `json:"seqNo"`
	Documents                    []interface{}   `json:"documents"`
	JoinMaxLifetimeLimitCatCode  string          `json:"joinMaxLifetimeLimitCatCode"`
	Configurations               []interface{}   `json:"configurations"`
	Points                       int             `json:"points"`
	ProductLine                  string          `json:"productLine"`
	JoinMaxLifetimeLimit         int             `json:"joinMaxLifetimeLimit"`
	AdditionalInfo               []interface{}   `json:"additionalInfo"`
	Sku                          string          `json:"sku"`
	ProductCMSData               []interface{}   `json:"productCmsData"`
	ActiveSmartDelivery          bool            `json:"activeSmartDelivery"`
	BrandName                    string          `json:"brandName"`
	DoNotSplitPackBV             bool            `json:"doNotSplitPackBV"`
	SDOnlyPackage                bool            `json:"sdOnlyPackage"`
	Weight                       int             `json:"weight"`
	IsVolumeBasedRSB             bool            `json:"isVolumeBasedRSB"`
	IsRetailCart                 bool            `json:"isRetailCart"`
	OfferLoyaltyProgram          bool            `json:"offerLoyaltyProgram"`
	IsComingSoon                 bool            `json:"isComingSoon"`
	BrandID                      int             `json:"brandId"`
	ProductMenuID                int             `json:"productMenuId"`
	Name                         string          `json:"name"`
	Messages                     []Message       `json:"messages"`
	IsRedemption                 bool            `json:"isRedemption"`
	Description                  string          `json:"description"`
	AutoshipProductPk            int             `json:"autoshipProductPk"`
	OfferPreferredCust           bool            `json:"offerPreferredCust"`
	ProductCategory              string          `json:"productCategory"`
	Bv                           int             `json:"bv"`
	ShowSDCheckbox               bool            `json:"showSDCheckbox"`
	OfferSDOnShop                bool            `json:"offerSDOnShop"`
	ImageURL                     string          `json:"imageUrl"`
	OfferAffiliateProgram        bool            `json:"offerAffiliateProgram"`
	SP                           int             `json:"sp"`
	IsShippable                  bool            `json:"isShippable"`
	ProductLineID                int             `json:"productLineId"`
	ProductPK                    int             `json:"productPK"`
	IsConfigurable               bool            `json:"isConfigurable"`
	IsStarterKit                 bool            `json:"isStarterKit"`
	IsRetailPackage              bool            `json:"isRetailPackage"`
	PriceType                    string          `json:"priceType"`
	ProductFunction              string          `json:"productFunction"`
	IsPackage                    bool            `json:"isPackage"`
	MainType                     int             `json:"mainType"`
	ProductCode                  string          `json:"productCode"`
	PackageItems                 []interface{}   `json:"packageItems"`
	IsFoodProduct                bool            `json:"isFoodProduct"`
	MaxLimit                     int             `json:"maxLimit"`
	ImageUrls                    []ImageURL      `json:"imageUrls"`
	IsProductAvailableOnAutoship bool            `json:"isProductAvailableOnAutoship"`
	ProductMenu                  string          `json:"productMenu"`
	Pricing                      []PricingModule `json:"pricing"`
}

type ImageURL struct {
	ImageName string `json:"imageName"`
	ImageURL  string `json:"imageUrl"`
}

type Message struct {
	TranslatedMessage string `json:"translatedMessage"`
	MessageType       string `json:"messageType"`
	MessageVariable   string `json:"messageVariable"`
}

type PricingModule struct {
	FormattedPrice string `json:"formattedPrice"`
	Price          int    `json:"price"`
	PriceType      string `json:"priceType"`
	PriceWarning   string `json:"priceWarning"`
	CurrencySymbol string `json:"currencySymbol"`
	NoVatPrice     int    `json:"noVatPrice"`
}
