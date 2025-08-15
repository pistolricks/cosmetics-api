package main

import (
	"net/http"
)

const shopifyAccessTokenHeader = "X-Shopify-Access-Token"

func (app *application) RoundTrip(req *http.Request) (*http.Response, error) {
	isAccessTokenSet := app.transport.accessToken != ""
	areBasicAuthCredentialsSet := app.transport.apiKey != "" && isAccessTokenSet

	if areBasicAuthCredentialsSet {
		req.SetBasicAuth(app.transport.apiKey, app.transport.accessToken)
	} else if isAccessTokenSet {
		req.Header.Set(shopifyAccessTokenHeader, app.transport.accessToken)
	}

	return app.transport.roundTripper.RoundTrip(req)
}
