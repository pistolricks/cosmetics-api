package main

import (
	"net/http"
	"os"
)

type Credentials struct {
	GrantType    string
	ClientId     string
	ClientSecret string
}

type CredentialsResponse struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	IssuedAt        int64  `json:"issued_at"`
	ExpiresIn       int    `json:"expires_in"`
	Status          string `json:"status"`
	Scope           string `json:"scope"`
	Issuer          string `json:"issuer"`
	ClientId        string `json:"client_id"`
	ApplicationName string `json:"application_name"`
	ApiProducts     string `json:"api_products"`
	PublicKey       string `json:"public_key"`
}

type Refresh struct {
	GrantType    string
	ClientId     string
	ClientSecret string
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	IssuedAt              int64  `json:"issued_at"`
	ExpiresIn             int    `json:"expires_in"`
	Status                string `json:"status"`
	Scope                 string `json:"scope"`
	Issuer                string `json:"issuer"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenIssuedAt  int64  `json:"refresh_token_issued_at"`
	RefreshTokenStatus    string `json:"refresh_token_status"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	RefreshCount          int    `json:"refresh_count"`
	ClientId              string `json:"client_id"`
	ApplicationName       string `json:"application_name"`
	ApiProducts           string `json:"api_products"`
	PublicKey             string `json:"public_key"`
}

type AuthorizationExchange struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type AuthorizationExchangeResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	IssuedAt              int64  `json:"issued_at"`
	ExpiresIn             int    `json:"expires_in"`
	Status                string `json:"status"`
	Scope                 string `json:"scope"`
	Issuer                string `json:"issuer"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenIssuedAt  int64  `json:"refresh_token_issued_at"`
	RefreshTokenStatus    string `json:"refresh_token_status"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	RefreshCount          int    `json:"refresh_count"`
	ClientId              string `json:"client_id"`
	ApplicationName       string `json:"application_name"`
	ApiProducts           string `json:"api_products"`
	PublicKey             string `json:"public_key"`
}

type AuthorizationCode struct {
	ClientId     string `json:"client_id"`
	ResponseType string `json:"response_type"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

type AuthorizationCodeResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type Revoke struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
}

type RevokeResponse struct {
	Error              string `json:"error"`
	ErrorDescription   string `json:"error_description"`
	ErrorUri           string `json:"error_uri"`
	RefreshTokenStatus string `json:"refresh_token_status"`
}

type CredentialsError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

type AddressResponse struct {
	Firm    string `json:"firm"`
	Address struct {
		StreetAddress             string `json:"streetAddress"`
		StreetAddressAbbreviation string `json:"streetAddressAbbreviation"`
		SecondaryAddress          string `json:"secondaryAddress"`
		CityAbbreviation          string `json:"cityAbbreviation"`
		City                      string `json:"city"`
		State                     string `json:"state"`
		ZIPCode                   string `json:"ZIPCode"`
		ZIPPlus4                  string `json:"ZIPPlus4"`
		Urbanization              string `json:"urbanization"`
	} `json:"address"`
	AdditionalInfo struct {
		DeliveryPoint        string `json:"deliveryPoint"`
		CarrierRoute         string `json:"carrierRoute"`
		DPVConfirmation      string `json:"DPVConfirmation"`
		DPVCMRA              string `json:"DPVCMRA"`
		Business             string `json:"business"`
		CentralDeliveryPoint string `json:"centralDeliveryPoint"`
		Vacant               string `json:"vacant"`
	} `json:"additionalInfo"`
	Corrections []struct {
		Code string `json:"code"`
		Text string `json:"text"`
	} `json:"corrections"`
	Matches []struct {
		Code string `json:"code"`
		Text string `json:"text"`
	} `json:"matches"`
	Warnings []string `json:"warnings"`
}

func (app *application) createToken(w http.ResponseWriter, r *http.Request) {

	credentials := Credentials{
		GrantType:    "client_credentials",
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(Credentials{
			GrantType:    credentials.GrantType,
			ClientId:     credentials.ClientId,
			ClientSecret: credentials.ClientSecret,
		}).SetResult(&CredentialsResponse{}). // or SetResult(LoginResponse{}).
		SetError(&CredentialsError{}).        // or SetError(LoginError{}).
		Post("https://apis.usps.com/oauth2/v3/token")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {

	credentials := Refresh{
		GrantType:    "refresh_token",
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(Refresh{
			GrantType:    credentials.GrantType,
			ClientId:     credentials.ClientId,
			ClientSecret: credentials.ClientSecret,
			RefreshToken: credentials.RefreshToken,
		}).SetResult(&RefreshResponse{}). // or SetResult(LoginResponse{}).
		SetError(&CredentialsError{}).    // or SetError(LoginError{}).
		Post("https://apis.usps.com/oauth2/v3/token")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) authorizationExchange(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Code        string `json:"code"`
		RedirectUri string `json:"redirect_uri"`
	}

	credentials := AuthorizationExchange{
		GrantType:    "authorization_code",
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Code:         input.Code,
		RedirectUri:  input.RedirectUri,
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(AuthorizationExchange{
			GrantType:    credentials.GrantType,
			ClientId:     credentials.ClientId,
			ClientSecret: credentials.ClientSecret,
			Code:         credentials.Code,
			RedirectUri:  credentials.RedirectUri,
		}).SetResult(&AuthorizationExchangeResponse{}). // or SetResult(LoginResponse{}).
		SetError(&CredentialsError{}).                  // or SetError(LoginError{}).
		Post("https://apis.usps.com/oauth2/v3/token")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createAuthorizationCode(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Code        string `json:"code"`
		RedirectUri string `json:"redirect_uri"`
		Scope       string `json:"scope"`
		State       string `json:"state"`
	}

	credentials := AuthorizationCode{
		ClientId:     os.Getenv("CLIENT_ID"),
		ResponseType: "code",
		RedirectUri:  input.RedirectUri,
		Scope:        input.Scope,
		State:        input.State,
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(AuthorizationCode{
			ClientId:     credentials.ClientId,
			ResponseType: credentials.ResponseType,
			RedirectUri:  credentials.RedirectUri,
			Scope:        credentials.Scope,
			State:        credentials.State,
		}).SetResult(&AuthorizationCodeResponse{}). // or SetResult(LoginResponse{}).
		SetError(&CredentialsError{}).              // or SetError(LoginError{}).
		Post("https://apis.usps.com/oauth2/v3/authorize")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) revokeToken(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token         string `json:"token"`
		TokenTypeHint string `json:"refresh_token"`
	}

	credentials := Revoke{
		Token:         input.Token,
		TokenTypeHint: input.TokenTypeHint,
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(Revoke{
			Token:         credentials.Token,
			TokenTypeHint: credentials.TokenTypeHint,
		}).SetResult(&RevokeResponse{}). // or SetResult(LoginResponse{}).
		SetError(&CredentialsError{}).   // or SetError(LoginError{}).
		Post("https://apis.usps.com/oauth2/v3/revoke")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) addressLookup(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token            string `json:"token"`
		Firm             string `json:"firm"`
		StreetAddress    string `json:"streetAddress"`
		SecondaryAddress string `json:"secondaryAddress"`
		City             string `json:"city"`
		State            string `json:"state"`
		Urbanization     string `json:"urbanization"`
		ZIPCode          string `json:"ZIPCode"`
		ZIPPlus4         string `json:"ZIPPlus4"`
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetQueryParams(map[string]string{
			"firm":             input.Firm,
			"streetAddress":    input.StreetAddress,
			"secondaryAddress": input.SecondaryAddress,
			"city":             input.City,
			"state":            input.State,
			"urbanization":     input.Urbanization,
			"zipCode":          input.ZIPCode,
			"zipPlus4":         input.ZIPPlus4,
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken(input.Token).
		SetResult(&AddressResponse{}).
		SetError(&CredentialsError{}).
		Get("https://apis.usps.com/addresses/v3/address")

	err = app.writeJSON(w, http.StatusOK, envelope{"resources": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
