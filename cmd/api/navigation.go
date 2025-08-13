package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/validator"
	"golang.org/x/net/context"
)

func (app *application) chromeLoginHandler(w http.ResponseWriter, r *http.Request) {

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

	page, browser, cookies := app.chromium.Chrome.ChromeLogin(input.LoginUrl, input.RimanStoreName, credentials.UserName, credentials.Password)

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

	err = app.writeJSON(w, http.StatusOK, envelope{"client": client, "page": page, "browser": browser, "cookies": cookies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) chromeHomePageHandler(w http.ResponseWriter, r *http.Request) {
	// networkCookie := networkCookies(cookies)

	rimanStoreName := app.envars.RimanStoreName // os.Getenv("RIMAN_STORE_NAME")
	rimanRid := app.envars.Username             // os.Getenv("USERNAME")

	cookies := app.chromium.Chrome.ChromeHomePage(rimanStoreName)

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

func (app *application) apiRimanLogoutHandler(w http.ResponseWriter, r *http.Request) {

	p := app.browser
	p.MustClose()

	/*
			res, err := riman.Logout(app.envars.Token)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}

		fmt.Println(res)
		fmt.Println(err)
	*/

	err := app.writeJSON(w, http.StatusOK, envelope{"status": "logged_out"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func networkCookies(cookies []*proto.NetworkCookie) []*proto.NetworkCookieParam {

	var networkCookie []*proto.NetworkCookieParam

	for _, cookie := range cookies {
		networkCookie = append(networkCookie, &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
			SameSite: cookie.SameSite,
			Expires:  cookie.Expires,
		})
	}

	return networkCookie
}

func (app *application) processBilling(page *rod.Page) {

	page.MustActivate()

}

type StateObject = struct {
	Code  string
	Name  string
	name2 any
}

/* TODO: REMOVE HARD CODED EMAIL */

func IsValidPhoneNumber(phone_number string) bool {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phone_number = strings.ReplaceAll(phone_number, " ", "")

	return re.Find([]byte(phone_number)) != nil
}

func (app *application) emulateGPS(page *rod.Page) bool {
	latStr := os.Getenv("LAT")
	lngStr := os.Getenv("LNG")

	// Convert string coordinates to float64
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
	}

	accuracy := 100.0 // Define accuracy in meters

	override := &proto.EmulationSetGeolocationOverride{
		Latitude:  &lat,
		Longitude: &lng,
		Accuracy:  &accuracy, // Accuracy is optional but recommended.
	}

	_, err = page.Call(context.Background(), "", "Emulation.setGeolocationOverride", override)
	if err != nil {
		fmt.Println("Failed to set geolocation:", err)
	}

	return true
}

func (app *application) submitFormHandler(w http.ResponseWriter, r *http.Request) {

	el := app.chromium.Chrome.SubmitForm()

	err := app.writeJSON(w, http.StatusOK, envelope{"element": el}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
