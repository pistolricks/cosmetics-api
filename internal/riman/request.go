package riman

import (
	"fmt"
	"net/url"

	"resty.dev/v3"
)

type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoggedInResponse struct {
	SecurityRedirect bool   `json:"securityRedirect"`
	Status           string `json:"-"`
	LiToken          string `json:"liToken"`
	LiUser           string `json:"liUser"`
	Jwt              string `json:"jwt"`
}

const loginUrl = "https://security-api.riman.com/api/v2/CheckAttemptsAndLogin"

func (m *ClientModel) Login(token string, credentials Credentials) (*LoggedInResponse, error) {

	params := url.Values{}
	params.Add("userName", credentials.UserName)
	params.Add("password", credentials.Password)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetBody(Credentials{
			UserName: credentials.UserName,
			Password: credentials.Password,
		}).
		SetResult(&LoggedInResponse{}).
		SetError(&Errors{}).
		Post(loginUrl)

	return res.Result().(*LoggedInResponse), err
}

type ReissueTokenResponse = map[string]any

func (m *ClientModel) ReissueToken(token string) (*ReissueTokenResponse, error) {

	logoutUrl := fmt.Sprintf("https://security-api.riman.com/api/v2/token/reissue")

	fmt.Println(logoutUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetResult(&ReissueTokenResponse{}).
		SetError(&Errors{}).
		Post(logoutUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*ReissueTokenResponse))

	return res.Result().(*ReissueTokenResponse), err

}

type LogoutResponse = map[string]any

func (m *ClientModel) Logout(token string) (*LogoutResponse, error) {

	logoutUrl := fmt.Sprintf("https://security-api.riman.com/api/v2/token/logout")

	fmt.Println(logoutUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetResult(&LogoutResponse{}).
		SetError(&Errors{}).
		Post(logoutUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*LogoutResponse))

	return res.Result().(*LogoutResponse), err
}
