package riman

type ShippedProductResponse struct {
	PackagePk                 int         `json:"packagePk"`
	ProductPk                 int         `json:"productPk"`
	PackageName               string      `json:"packageName"`
	ProductName               string      `json:"productName"`
	IsPackage                 bool        `json:"isPackage"`
	Quantity                  int         `json:"quantity"`
	Cv                        int         `json:"cv"`
	Sp                        int         `json:"sp"`
	Price                     int         `json:"price"`
	FormattedPrice            string      `json:"formattedPrice"`
	CurrencyCode              string      `json:"currencyCode"`
	ShipmentNumber            string      `json:"shipmentNumber"`
	ShipmentStatus            string      `json:"shipmentStatus"`
	ShippedDate               string      `json:"shippedDate"`
	TrackingNumber            string      `json:"trackingNumber"`
	TrackingLink              string      `json:"trackingLink"`
	VideoOrderPackagingInfoPK interface{} `json:"videoOrderPackagingInfoPK"`
}

type PackageItem struct {
	ProductPK   int    `json:"productPK"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImgUrl      string `json:"imgUrl"`
	Qty         int    `json:"qty"`
}

type Pricing struct {
	PriceType      string `json:"priceType"`
	CurrencySymbol string `json:"currencySymbol"`
	Price          int    `json:"price"`
	NoVatPrice     int    `json:"noVatPrice"`
	FormattedPrice string `json:"formattedPrice"`
	PriceWarning   string `json:"priceWarning"`
}

type ImageUrl struct {
	ImageUrl  string `json:"imageUrl"`
	ImageName string `json:"imageName"`
}

type ProductCmsData struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	DataTag string `json:"dataTag"`
}

type RimanProduct struct {
	ProductPK                    int              `json:"productPK"`
	ProductCode                  string           `json:"productCode"`
	Sku                          interface{}      `json:"sku"`
	ProductCategory              string           `json:"productCategory"`
	BrandId                      int              `json:"brandId"`
	BrandName                    string           `json:"brandName"`
	ProductBrand                 interface{}      `json:"productBrand"`
	Name                         string           `json:"name"`
	ImageUrl                     string           `json:"imageUrl"`
	Weight                       int              `json:"weight"`
	IsComingSoon                 bool             `json:"isComingSoon"`
	ComingSoonMessage            interface{}      `json:"comingSoonMessage"`
	IsPackage                    bool             `json:"isPackage"`
	PackageItems                 []PackageItem    `json:"packageItems"`
	IsConfigurable               bool             `json:"isConfigurable"`
	IsProductAvailableOnAutoship bool             `json:"isProductAvailableOnAutoship"`
	IsProductOnAutoship          bool             `json:"isProductOnAutoship"`
	AutoshipProductPk            int              `json:"autoshipProductPk"`
	MaxLimit                     int              `json:"maxLimit"`
	Points                       float64          `json:"points"`
	Bv                           int              `json:"bv"`
	Sp                           int              `json:"sp"`
	ProductMenuId                int              `json:"productMenuId"`
	ProductMenu                  string           `json:"productMenu"`
	Configurations               []interface{}    `json:"configurations"`
	Pricing                      []Pricing        `json:"pricing"`
	Description                  string           `json:"description"`
	ImageUrls                    []ImageUrl       `json:"imageUrls"`
	AdditionalInfo               []interface{}    `json:"additionalInfo"`
	Documents                    []interface{}    `json:"documents"`
	IsShippable                  bool             `json:"isShippable"`
	IsStarterKit                 bool             `json:"isStarterKit"`
	SeqNo                        int              `json:"seqNo"`
	IsFoodProduct                bool             `json:"isFoodProduct"`
	RankInfo                     interface{}      `json:"rankInfo"`
	IsRetailPackage              bool             `json:"isRetailPackage"`
	IsVolumeBasedRSB             bool             `json:"isVolumeBasedRSB"`
	MainType                     int              `json:"mainType"`
	ActiveSmartDelivery          bool             `json:"activeSmartDelivery"`
	IsRedemption                 bool             `json:"isRedemption"`
	PriceType                    string           `json:"priceType"`
	DoNotSplitPackBV             bool             `json:"doNotSplitPackBV"`
	SdOnlyPackage                bool             `json:"sdOnlyPackage"`
	ShowSDCheckbox               bool             `json:"showSDCheckbox"`
	OfferAffiliateProgram        bool             `json:"offerAffiliateProgram"`
	OfferPreferredCust           bool             `json:"offerPreferredCust"`
	OfferLoyaltyProgram          bool             `json:"offerLoyaltyProgram"`
	OfferSDOnShop                bool             `json:"offerSDOnShop"`
	IsRetailCart                 bool             `json:"isRetailCart"`
	ProductLineId                int              `json:"productLineId"`
	ProductLine                  string           `json:"productLine"`
	ProductFunction              string           `json:"productFunction"`
	MaxLifetimeLimitCatCode      interface{}      `json:"maxLifetimeLimitCatCode"`
	MaxLifetimeLimit             interface{}      `json:"maxLifetimeLimit"`
	JoinMaxLifetimeLimitCatCode  string           `json:"joinMaxLifetimeLimitCatCode"`
	JoinMaxLifetimeLimit         int              `json:"joinMaxLifetimeLimit"`
	ProductCmsData               []ProductCmsData `json:"productCmsData"`
}

type Order struct {
	OrderDate               string  `json:"orderDate"`
	MainOrdersPK            int     `json:"mainOrdersPK"`
	OrderType               string  `json:"orderType"`
	FinalOrderType          string  `json:"finalOrderType,omitempty"`
	SiteURL                 string  `json:"siteURL"`
	EncOrderNumber          string  `json:"encOrderNumber"`
	CurrencySymbol          string  `json:"currencySymbol"`
	CurrencyCode            string  `json:"currencyCode"`
	PaidStatus              bool    `json:"paidStatus"`
	HasTaxInvoice           bool    `json:"hasTaxInvoice"`
	HasCommercialInvoice    bool    `json:"hasCommercialInvoice"`
	HasCreditNote           bool    `json:"hasCreditNote"`
	IsShippingPending       bool    `json:"isShippingPending"`
	IsPB                    bool    `json:"isPB"`
	IsPA                    bool    `json:"isPA"`
	IsCC                    bool    `json:"isCC"`
	MainFK                  int     `json:"mainFK"`
	MainOrderTypeFK         int     `json:"mainOrderTypeFK"`
	VoucherURL              string  `json:"voucherURL,omitempty"`
	ShipCountry             string  `json:"shipCountry"`
	ShippingStatus          string  `json:"shippingStatus"`
	OrderShippingStatus     string  `json:"orderShippingStatus"`
	OrderTypeValue          string  `json:"orderTypeValue,omitempty"`
	PaidStatusValue         string  `json:"paidStatusValue"`
	Quantity                int     `json:"quantity"`
	Email                   string  `json:"email,omitempty"`
	Phone                   string  `json:"phone,omitempty"`
	ShipFirstName           string  `json:"shipFirstName,omitempty"`
	ShipLastName            string  `json:"shipLastName,omitempty"`
	MarkedPaidDate          string  `json:"markedPaidDate"`
	Total                   float64 `json:"total"`
	ConvTotal               float64 `json:"convTotal"`
	ConvTotalFormat         string  `json:"convTotalFormat"`
	SubTotal                float64 `json:"subTotal"`
	ConvSubtotal            float64 `json:"convSubtotal"`
	ShipTotal               float64 `json:"shipTotal"`
	ConvShipTotal           float64 `json:"convShipTotal"`
	Taxes                   float64 `json:"taxes"`
	TaxLabel                string  `json:"taxLabel"`
	ProductTax              float64 `json:"productTax"`
	ShippingTax             float64 `json:"shippingTax"`
	TotalProductTax         float64 `json:"totalProductTax"`
	AdditionalTaxLabel      string  `json:"additionalTaxLabel"`
	AdditionalTax           string  `json:"additionalTax,omitempty"`
	ConvTaxes               float64 `json:"convTaxes"`
	OrderProcessingFees     string  `json:"orderProcessingFees,omitempty"`
	ConvOrderProcessingFees string  `json:"convOrderProcessingFees,omitempty"`
	Discount                float64 `json:"discount"`
	ConvDiscount            float64 `json:"convDiscount"`
	RefundAmount            float64 `json:"refundAmount"`
	ConvRefund              float64 `json:"convRefund"`
	SalesCampaignFK         int     `json:"salesCampaignFK,omitempty"`
	Paidstatusfk            int     `json:"paidstatusfk"`
	DeliveryDate            string  `json:"deliveryDate,omitempty"`
	ShippingDetails         string  `json:"shippingDetails,omitempty"`
	OrderItems              string  `json:"orderItems,omitempty"`
	Payments                string  `json:"payments,omitempty"`
	IsPrepaidOrder          string  `json:"isPrepaidOrder,omitempty"`
	ShowInvoice             bool    `json:"showInvoice"`
	ShowOrderInvoice        bool    `json:"showOrderInvoice"`
	KrGuaranteeNo           string  `json:"krGuaranteeNo"`
	WeChatOrderNumber       string  `json:"weChatOrderNumber,omitempty"`
	MemberID                string  `json:"memberID,omitempty"`
}
