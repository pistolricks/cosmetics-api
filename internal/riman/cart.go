package riman

import (
	"database/sql"
	"fmt"

	"resty.dev/v3"
)

type CartErrors struct {
	Error string `json:"error"`
}
type Body map[string]any

type CartModel struct {
	DB *sql.DB
}

func GetCart(cartKey string) (*Cart, error) {

	cartUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s", cartKey)

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

	fmt.Println(err)

	return res.Result().(*Cart), err
}

func AddProductToCart(token string, cartKey string, addProductPayload *AddProductPayload) (*CartItem, error) {

	addProductUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s/items", cartKey)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(Body{
			"configFk":        addProductPayload.ConfigFk,
			"discount":        addProductPayload.Discount,
			"extraFee":        addProductPayload.ExtraFee,
			"mainCartFk":      addProductPayload.MainCartFk,
			"mainCartItemsPk": addProductPayload.MainCartItemsPk,
			"productFk":       addProductPayload.ProductFk,
			"quantity":        addProductPayload.Quantity,
			"setupForAs":      addProductPayload.SetupForAs,
		}).
		SetResult(&CartItem{}).
		Post(addProductUrl)

	fmt.Println(err, res)
	cartResponse := res.Result().(*CartItem)

	return cartResponse, err
}

func PatchLocale(token string, cartKey string, mainFk int) (*Cart, error) {

	cartUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s", cartKey)

	fmt.Println(cartUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetBody(&Body{
			"countryCode": "US",
			"culture":     "en-US",
			"mainFk":      mainFk,
		}).
		SetResult(&Cart{}).
		SetError(&CartErrors{}).
		Patch(cartUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*Cart))

	return res.Result().(*Cart), err
}

func Patch(token string, cartKey string) (*Cart, error) {

	cartUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s", cartKey)

	fmt.Println(cartUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetBody(&Body{
			"shippingTypeFk": 1163,
		}).
		SetResult(&Cart{}).
		SetError(&CartErrors{}).
		Patch(cartUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*Cart))

	return res.Result().(*Cart), err
}

func DeleteProductFromCart(token string, cartKey string, id string) (string, error) {

	deleteProductUrl := fmt.Sprintf("https://cart-api.riman.com/api/v1/shopping/%s/items/%s", cartKey, id)

	fmt.Println(deleteProductUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		Delete(deleteProductUrl)

	fmt.Println(err, res)

	var d = "Deleted"

	return d, err
}

func AddProductToCartWithQuantity(token string, cartKey string, payload *AddProductPayload) (*CartItem, error) {
	ci, err := AddProductToCart(token, cartKey, payload)
	if err != nil {
		return ci, err
	}

	return ci, err
}
