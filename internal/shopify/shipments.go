package shopify

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

type Address struct {
	RegionCode   string   `json:"regionCode"`
	Locality     string   `json:"locality"`
	AddressLines []string `json:"addressLines"`
}

type Body map[string]any
type Results map[string]any
type Errors = struct {
	errors error
}

func AddressValidation(address *Address) (*AddressValidationResponse, error) {

	fmt.Println("ADDRESS --")
	fmt.Println(address)

	client := resty.New()
	defer client.Close()

	apiKey := os.Getenv("GMAPS_KEY")

	fmt.Println(apiKey)
	addressUrl := fmt.Sprintf("https://addressvalidation.googleapis.com/v1:validateAddress?key=%s", apiKey)

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(Body{
			"regionCode":   address.RegionCode,
			"locality":     address.Locality,
			"addressLines": address.AddressLines,
		}).
		SetResult(&AddressValidationResponse{}). // or SetResult(LoginResponse{}).
		SetError(&Errors{}).                     // or SetError(LoginError{}).
		Post(addressUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.Result().(*AddressValidationResponse))

	return res.Result().(*AddressValidationResponse), err
}
