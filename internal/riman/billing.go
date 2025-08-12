package riman

type BillingAddress struct {
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
