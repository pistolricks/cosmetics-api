package data

type Order struct {
	MainId          int         `json:"mainId"`
	MainOrderType   int         `json:"mainOrderType"`
	CountryCode     string      `json:"countryCode"`
	SalesCampaignFK interface{} `json:"salesCampaignFK"`
	CartKey         string      `json:"cartKey"`
}

type OrderResponse struct {
	MainOrdersFK   int    `json:"mainOrdersFK"`
	EncOrderNumber string `json:"encOrderNumber"`
	Message        string `json:"message"`
}

type Errors struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

func OrderCreate(order Order) (*OrderResponse, error) {

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(Order{
			MainId:          order.MainId,
			MainOrderType:   order.MainOrderType,
			CountryCode:     order.CountryCode,
			SalesCampaignFK: order.SalesCampaignFK,
			CartKey:         order.CartKey,
		}).                          // default request content type is JSON
		SetResult(&OrderResponse{}). // or SetResult(LoginResponse{}).
		SetError(&Errors{}).         // or SetError(LoginError{}).
		Post("https://cart-api.riman.com/api/v2/order")

	if err != nil {
		return nil, err
	}

	return res.Result().(*OrderResponse), err
}
