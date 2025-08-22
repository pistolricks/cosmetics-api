package riman

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type RimanBillingAddress struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Phone     string `json:"phone"`
}

type RimanCreditCard struct {
	CardName   string `json:"cardName"`
	CardNumber string `json:"cardNumber"`
	ExpMonth   string `json:"expMonth"`
	ExpYear    string `json:"expYear"`
	CVV        string `json:"cvv"`
}

type State struct {
	Code  string      `json:"code"`
	Name  string      `json:"name"`
	Name2 interface{} `json:"-"`
}

func ValidateBilling() {

}
