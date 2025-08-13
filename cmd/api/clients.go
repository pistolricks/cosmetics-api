package main

import (
	"fmt"
	"net/http"
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
