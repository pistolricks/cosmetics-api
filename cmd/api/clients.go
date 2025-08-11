package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/validator"
)

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

func (app *application) findCookieValue() *string {
	for i := range app.cookies {
		if app.cookies[i].Name == "token" {
			app.envars.Token = app.cookies[i].Value
			fmt.Println("app.envars.Token")
			fmt.Println(app.envars.Token)
			return &app.cookies[i].Value
		}
	}
	// Return nil if no product is found
	return nil
}

func (app *application) findCartKeyValue() *string {
	for i := range app.cookies {
		if app.cookies[i].Name == "cartKey" {
			return &app.cookies[i].Value
		}
	}
	// Return nil if no product is found
	return nil
}

func (app *application) clientLoginHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		RimanStoreName string `json:"rimanStoreName"`
		UserName       string `json:"userName"`
		Password       string `json:"password"`
		LoginUrl       string `json:"loginUrl"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	credentials := riman.Credentials{
		UserName: input.UserName,
		Password: input.Password,
	}

	v := validator.New()
	data.ValidatePasswordPlaintext(v, credentials.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.envars.RimanStoreName = input.RimanStoreName
	app.envars.Username = credentials.UserName
	app.envars.Password = credentials.Password
	app.envars.LoginUrl = input.LoginUrl

	// err = os.Setenv("RIMAN_STORE_NAME", input.RimanStoreName)
	// err = os.Setenv("LOGIN_URL", input.LoginUrl)
	// err = os.Setenv("USERNAME", credentials.UserName)
	// err = os.Setenv("Password", credentials.Password)

	page, browser, cookies := app.RimanLogin(input.LoginUrl, input.RimanStoreName, credentials.UserName, credentials.Password)

	app.page = page
	app.browser = browser
	app.cookies = cookies
	fmt.Println(browser)

	client, err := app.riman.Clients.GetByClientUsername(input.UserName)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.client = client

	err = app.writeJSON(w, http.StatusOK, envelope{"client": client, "page": app.page, "browser": app.browser, "cookies": cookies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {

	rimanStoreName := app.envars.RimanStoreName // os.Getenv("RIMAN_STORE_NAME")
	rimanRid := app.envars.Username             // os.Getenv("USERNAME")
	currentPage := app.page
	currentBrowser := app.browser

	_, _, cookies := app.HomePage(rimanStoreName, currentPage, currentBrowser)

	fmt.Println(cookies)

	app.cookies = cookies

	token := app.findCookieValue()
	if token == nil {
		app.invalidCredentialsResponse(w, r)
		return
	}
	fmt.Println("TOKEN")
	fmt.Println(token)

	/* ADD SESSION HERE */

	client, err := app.riman.Clients.GetByClientUsername(rimanRid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.client = client

	cookie := app.findCookieValue()
	if cookie == nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	cartKey := app.findCartKeyValue()
	if cartKey == nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	session, err := app.riman.Session.NewRimanSession(client.ID, 24*time.Hour, riman.ScopeAuthentication, app.envars.Token, *cartKey, envelope{"cookies": cookies})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.session = session

	fmt.Println("SESSION CLIENT ID")
	fmt.Println(session.ClientID)

	err = app.writeJSON(w, http.StatusOK, envelope{"session": session, "client": client, "errors": err}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClientsHandler(w http.ResponseWriter, r *http.Request) {

	clients, metadata, err := app.riman.Clients.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"clients": clients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
